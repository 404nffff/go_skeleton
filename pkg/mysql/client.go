package mysql

import (
	"fmt"
	"log"
	"sync"
	"time"
	"tool/pkg/event_manage"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局 sync.Map 变量
var (
	clients sync.Map
)

// NewClient 初始化 GORM 客户端，并支持多个数据库连接
func NewClient(name string) *gorm.DB {
	db, loaded := clients.LoadOrStore(name, createDBClient(name))
	if loaded {
		// 检查连接是否有效
		if !isValidConnection(db.(*gorm.DB)) {
			log.Printf("数据库连接丢失，正在重新连接: %s", name)
			db = createDBClient(name)
			clients.Store(name, db)
		} else {
			printConnectionPoolStats(db.(*gorm.DB))
		}
	} else {
		printConnectionPoolStats(db.(*gorm.DB))
	}
	return db.(*gorm.DB)
}

// isValidConnection 检查数据库连接是否有效
func isValidConnection(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

// printConnectionPoolStats 打印连接池状态
func printConnectionPoolStats(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取数据库实例失败: %v", err)
		return
	}
	stats := sqlDB.Stats()
	log.Printf("数据库连接池状态: 最大连接数: %d, 打开连接数: %d, 空闲连接数: %d, 等待中的连接数: %d, 总连接数: %d, 最大生存时间: %d",
		stats.MaxOpenConnections, stats.OpenConnections, stats.Idle, stats.WaitCount, stats.MaxLifetimeClosed, stats.MaxIdleClosed)
}

// createDBClient 创建新的数据库客户端
func createDBClient(name string) *gorm.DB {

	// 加载配置
	config := loadConfig(name)

	// 判断是否为空
	if config.Host == "" {
		panic(fmt.Sprintf("获取 Mysql 配置失败: %s", name))
	}

	// 构建 DSN (数据源名称)
	dsn := config.User + ":" +
		config.Pass + "@tcp(" +
		config.Host + ":" +
		config.Port + ")/" +
		config.Database + "?charset=" +
		config.Charset + "&parseTime=True&loc=Local"

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(),
	})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 获取底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	// 设置数据库连接池配置
	sqlDB.SetMaxIdleConns(config.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.SetConnMaxLifetime) * time.Second)

	log.Printf("成功连接到数据库，数据库名称: %s", config.Database)

	// 配置 GORM 回调
	db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", func(d *gorm.DB) {
		d.Statement.RaiseErrorOnNotFound = false
	})

	// 创建事件管理工厂并注册销毁事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	eventName := config.EventDestroyPrefix
	if _, exists := eventManageFactory.Get(eventName); !exists {
		eventManageFactory.Set(eventName, func(args ...interface{}) {
			if err := sqlDB.Close(); err != nil {
				log.Printf("关闭 Mysql 连接失败: %v", err)
				return
			}
			log.Printf("销毁 Mysql 连接: %s", name)
		})
	}

	return db
}

// GetDB 获取数据库连接实例
func GetDB(name string) *gorm.DB {
	db, exists := clients.Load(name)
	if exists {
		return db.(*gorm.DB)
	}
	return nil
}
