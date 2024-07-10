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
var clients sync.Map

// DatabaseConfig 定义数据库配置结构体
type DatabaseConfig struct {
	User               string // 数据库用户名
	Pass               string // 数据库密码
	Host               string // 数据库地址
	Port               string // 数据库端口
	Database           string // 数据库名称
	Charset            string // 数据库字符集
	SetMaxIdleConns    int    // 连接池中的最大空闲连接数
	SetMaxOpenConns    int    // 数据库的最大连接数量
	SetConnMaxLifetime int    // 连接的最大可复用时间
	EventDestroyPrefix string // 事件销毁前缀
}

// NewDBClient 初始化 GORM 客户端，并支持多个数据库连接
func NewDBClient(name string, config DatabaseConfig) (*gorm.DB, error) {
	// 使用 LoadOrStore 确保在并发环境中只初始化一次数据库连接
	db, loaded := clients.LoadOrStore(name, createDBClient(name, config))
	if loaded {
		return db.(*gorm.DB), nil
	}

	return db.(*gorm.DB), nil
}

// createDBClient 创建新的数据库客户端
func createDBClient(name string, config DatabaseConfig) *gorm.DB {
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
	eventName := config.EventDestroyPrefix + "Mysql_" + name
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
