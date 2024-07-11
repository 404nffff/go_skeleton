package redis

import (
	"log"
	"tool/app/global/variable"
	"tool/pkg/yml_config"
)

type RedisConfig struct {
	Host                  string // Redis 服务器地址，格式为 "host:port"。
	Auth                  string // 可选的密码。如果不需要密码认证，请留空。
	IndexDb               int    // 数据库编号。默认为 0。
	PoolSize              int    // 每个 CPU 的最大连接数。默认为 10。
	MinIdleConns          int    // 最小空闲连接数。在建立新连接较慢时很有用。
	ConnFailRetryTimes    int    // 放弃前的最大重试次数。默认为不重试。
	ConnFailRetryInterval int    // 重试之间的最小退避时间。默认为 8 毫秒；-1 禁用退避。
	EventDestroyPrefix    string // 事件销毁前缀
}

func loadConfig(conn string) RedisConfig {

	// 加载配置文件
	redisConfig := yml_config.LoadConfig("redis")

	// 查找配置文件中的 Redis 配置
	// Local:
	// 	Host: "127.0.0.1:6311"
	// 	Auth: ""
	// 	IndexDb: 0
	// 	ConnFailRetryTimes: 1    #连接失败重试次数
	// 	ConnFailRetryInterval: 2 #连接失败重试间隔秒数
	// 	PoolSize: 5             #连接池大小
	// 	MinIdleConns: 2          #最小空闲连接数

	if redisConfig.GetString(conn+".Host") == "" {
		log.Fatalf("Failed to get Redis config: %s", conn)
	}

	config := RedisConfig{
		Host:                  redisConfig.GetString(conn + ".Host"),
		Auth:                  redisConfig.GetString(conn + ".Auth"),
		IndexDb:               redisConfig.GetInt(conn + ".IndexDb"),
		PoolSize:              redisConfig.GetInt(conn + ".PoolSize"),
		MinIdleConns:          redisConfig.GetInt(conn + ".MinIdleConns"),
		ConnFailRetryTimes:    redisConfig.GetInt(conn + ".ConnFailRetryTimes"),
		ConnFailRetryInterval: redisConfig.GetInt(conn + ".ConnFailRetryInterval"),
		EventDestroyPrefix:    variable.EventDestroyPrefix + "Redis_" + conn,
	}

	return config
}
