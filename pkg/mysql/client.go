package mysql

import (
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
	// 使用 LoadOrStore 确保在并发环境中只初始化一次数据库连接
	db, loaded := clients.LoadOrStore(name, createDBClient(name))
	if loaded {
		return db.(*gorm.DB)
	}

	return db.(*gorm.DB)
}

// createDBClient 创建新的数据库客户端
func createDBClient(name string) *gorm.DB {

	//加载配置
	config := loadConfig(name)

	//判断是否为空
	if config.Host == "" {
		log.Fatalf("Failed to get Mysql config: %s", name)
		return nil
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
		log.Fatalf("Failed to connect to database: %v", err)
		return nil
	}

	// 获取底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
		return nil
	}

	// 设置数据库连接池配置
	sqlDB.SetMaxIdleConns(config.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.SetConnMaxLifetime) * time.Second)

	log.Printf("Connected to database successfully, database: %s", config.Database)

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
				log.Printf("Failed to close Mysql connection: %v", err)
				return
			}
			log.Printf("Destroying Mysql connection for %s", name)
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
