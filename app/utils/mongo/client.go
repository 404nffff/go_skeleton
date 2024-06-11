package mongo

import (
	"context"
	"log"
	"sync"
	"time"
	"tool/app/utils/event_manage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 全局 sync.Map 变量
var clients sync.Map
var dbs sync.Map

// DatabaseConfig 定义数据库配置结构体
type DatabaseConfig struct {
	URI                string // 数据库连接 URI 字符串 (e.g. "mongodb://localhost:27017/")
	Database           string // 数据库名称
	MaxPoolSize        uint64 // 连接池中的最大连接数
	MinPoolSize        uint64 // 连接池中的最小连接数
	EventDestroyPrefix string // 事件销毁前缀
}

// InitMongo 初始化 MongoDB 客户端，并支持多个数据库连接
func InitMongo(dbConfig DatabaseConfig) (*mongo.Database, error) {
	// 使用 LoadOrStore 确保在并发环境中只初始化一次数据库连接
	db, loaded := dbs.LoadOrStore(dbConfig.Database, createMongoClient(dbConfig))
	if loaded {
		return db.(*mongo.Database), nil
	}

	return db.(*mongo.Database), nil
}

// createMongoClient 创建新的 MongoDB 客户端
func createMongoClient(dbConfig DatabaseConfig) *mongo.Database {
	// 设置客户端连接选项
	clientOptions := options.Client().
		ApplyURI(dbConfig.URI + dbConfig.Database).
		SetMaxPoolSize(dbConfig.MaxPoolSize).
		SetMinPoolSize(dbConfig.MinPoolSize).
		SetMonitor(NewMonitor())

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil
	}

	// 检查连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil
	}

	log.Printf("Connected to MongoDB successfully, database: %s", dbConfig.Database)

	// 注册销毁事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	eventName := dbConfig.EventDestroyPrefix + "Mongo_" + dbConfig.Database
	if _, exists := eventManageFactory.Get(eventName); !exists {
		eventManageFactory.Set(eventName, func(args ...interface{}) {
			CloseMongo(client, dbConfig.Database)
			log.Printf("Destroying MongoDB connection for %s", dbConfig.Database)
		})
	}

	clients.Store(dbConfig.Database, client)
	return client.Database(dbConfig.Database)
}

// GetCollection 获取指定数据库的集合
func GetCollection(dbName string, collection string) *mongo.Collection {
	db, exists := dbs.Load(dbName)
	if exists {
		return db.(*mongo.Database).Collection(collection)
	}
	return nil
}

// FindOne 查找单个文档
func FindOne(dbName string, collection string, filter interface{}) *mongo.SingleResult {
	col := GetCollection(dbName, collection)
	if col == nil {
		return nil
	}
	return col.FindOne(context.TODO(), filter)
}

// CloseMongo 关闭 MongoDB 客户端
func CloseMongo(client *mongo.Client, name string) {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
		log.Printf("Disconnected from MongoDB successfully, client: %s", name)
	}
}
