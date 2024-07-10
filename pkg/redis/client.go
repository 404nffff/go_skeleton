package redis

import (
	"context"
	"log"
	"sync"
	"time"
	"tool/pkg/event_manage"

	"github.com/go-redis/redis/v8"
)

// RedisClient 是一个全局的 Redis 客户端
var (
	clients sync.Map
	mu      sync.Mutex
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

func createClient(name string, config RedisConfig) *redis.Client {
	mu.Lock()
	defer mu.Unlock()

	options := &redis.Options{
		Network: "tcp",       // 网络类型，默认是 "tcp"，也可以使用 "unix"。
		Addr:    config.Host, // Redis 服务器地址，格式为 "host:port"。
		// 可选的用户名，用于 Redis ACL 支持。如果不使用 ACL，可以留空。
		Username: "",
		// 可选的密码。如果 Redis 服务器设置了密码认证，请在这里填写密码。
		// 如果没有设置密码认证，可以留空。
		Password: config.Auth,
		// 数据库编号，默认是 0。可以选择不同的数据库编号。
		DB: config.IndexDb,
		// 在放弃之前最大的重试次数。默认情况下，不重试失败的命令。
		MaxRetries: 3,
		// 每次重试之间的最小退避时间。默认是 8 毫秒；-1 禁用退避。
		MinRetryBackoff: 8 * time.Millisecond,
		// 每次重试之间的最大退避时间。默认是 512 毫秒；-1 禁用退避。
		MaxRetryBackoff: 512 * time.Millisecond,
		// 拨号超时时间，用于建立新连接。默认是 5 秒。
		DialTimeout: 5 * time.Second,
		// 读操作的超时时间。如果达到该时间，命令将失败并返回超时错误。
		// 默认是 3 秒。
		ReadTimeout: 3 * time.Second,
		// 写操作的超时时间。如果达到该时间，命令将失败并返回超时错误。
		// 默认是 ReadTimeout。
		WriteTimeout: 3 * time.Second,
		// 最大连接数，默认是每个 CPU 10 个连接。
		PoolSize: config.PoolSize,
		// 最小空闲连接数。保持这些空闲连接是有用的，特别是当建立新连接很慢时。
		MinIdleConns: config.MinIdleConns,
		// 连接的最大生存时间。到达此时间后，客户端将关闭连接。默认是不会关闭老连接。
		MaxConnAge: 0,
		// 空闲连接的超时时间。应该小于服务器的超时时间。默认是 5 分钟；-1 禁用空闲连接检查。
		IdleTimeout: 5 * time.Minute,
		// 空闲检查的频率。默认是 1 分钟；-1 禁用空闲连接检查。
		IdleCheckFrequency: 1 * time.Minute,
	}

	for i := 0; i < config.ConnFailRetryTimes; i++ {
		client := redis.NewClient(options)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := client.Ping(ctx).Err()
		if err == nil {
			log.Printf("Successfully connected to Redis")

			eventManageFactory := event_manage.CreateEventManageFactory()
			if _, exists := eventManageFactory.Get(config.EventDestroyPrefix + "Redis_" + name); !exists {
				eventManageFactory.Set(config.EventDestroyPrefix+"Redis_"+name, func(args ...interface{}) {
					_ = client.Close()
					log.Printf("Destroying Redis connection")
				})
			}

			return client
		}

		log.Printf("Failed to connect to Redis, retrying... (attempt %d)", i+1)
		time.Sleep(time.Duration(config.ConnFailRetryInterval) * time.Second)
	}

	log.Fatalf("Failed to connect to Redis, reached maximum retry attempts")
	return nil
}

func NewClient(name string, config RedisConfig) *redis.Client {
	client, _ := clients.LoadOrStore(name, createClient(name, config))
	return client.(*redis.Client)
}
