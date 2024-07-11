package memcached

import (
	"log"
	"tool/app/global/variable"
	"tool/pkg/yml_config"
)

type MemcachedConfig struct {
	Host                  string // Memcached 服务器地址，格式为 "host:port"。
	ConnFailRetryTimes    int    // 连接失败重试次数
	ConnFailRetryInterval int    // 连接失败重试间隔秒数
	EventDestroyPrefix    string // 事件销毁前缀
}

// 加载配置文件
func loadConfig(conn string) MemcachedConfig {

	// 加载配置文件
	memcachedConfig := yml_config.LoadConfig("memcached")

	// 查找配置文件中的 Memcached 配置
	// Local:
	// 	Host: "127.0.0.1:11213"
	// 	ConnFailRetryTimes: 1    #连接失败重试次数
	// 	ConnFailRetryInterval: 2 #连接失败重试间隔秒数

	if memcachedConfig.GetString(conn+".Host") == "" {
		log.Fatalf("Failed to get Memcached config: %s", conn)
	}

	config := MemcachedConfig{
		Host:                  memcachedConfig.GetString(conn + ".Host"),
		ConnFailRetryTimes:    memcachedConfig.GetInt(conn + ".ConnFailRetryTimes"),
		ConnFailRetryInterval: memcachedConfig.GetInt(conn + ".ConnFailRetryInterval"),
		EventDestroyPrefix:    variable.EventDestroyPrefix + "Memcached_" + conn,
	}

	return config
}
