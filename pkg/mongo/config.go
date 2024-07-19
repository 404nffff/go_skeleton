package mongo

import (
	"fmt"
	"tool/global/variable"
	"tool/pkg/yml_config"
)

// DatabaseConfig 定义数据库配置结构体
type DatabaseConfig struct {
	URI                string // 数据库连接 URI 字符串 (e.g. "mongodb://localhost:27017/")
	Database           string // 数据库名称
	MaxPoolSize        uint64 // 连接池中的最大连接数
	MinPoolSize        uint64 // 连接池中的最小连接数
	EventDestroyPrefix string // 事件销毁前缀
}

// 加载配置文件
func loadConfig(conn string) DatabaseConfig {

	mongoConfig := yml_config.LoadConfig("mongo")

	// 查找配置文件中的 MongoDB 配置
	// Local:
	// 	Open: true  # 是否启用 MongoDB
	// 	Uri: ""
	// 	Database: "hospital"
	// 	MaxPoolSize: 10  # 最大连接数
	// 	MinPoolSize: 1   # 最小空闲连接数

	if !mongoConfig.GetBool(conn + ".Open") {
		panic(fmt.Sprintf("Failed to get MongoDB config: %s", conn))
	}

	config := DatabaseConfig{
		URI:                mongoConfig.GetString(conn + ".Uri"),
		Database:           mongoConfig.GetString(conn + ".Database"),
		MaxPoolSize:        uint64(mongoConfig.GetInt(conn + ".MaxPoolSize")),
		MinPoolSize:        uint64(mongoConfig.GetInt(conn + ".MinPoolSize")),
		EventDestroyPrefix: variable.EventDestroyPrefix + "Mongo_" + conn,
	}

	return config
}
