package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"tool/pkg/event_manage"

	"github.com/go-redis/redis/v8"
)

// RedisClient 是一个全局的 Redis 客户端
var (
	clients sync.Map
)

// createClient 创建 Redis 客户端
func createClient(name string) *redis.Client {

	// 加载配置
	config := loadConfig(name)

	// 判断是否为空
	if config.Host == "" {
		panic(fmt.Sprintf("Failed to get Redis config: %s", name))
	}

	options := &redis.Options{
		Network:            "tcp",
		Addr:               config.Host,
		Username:           "",
		Password:           config.Auth,
		DB:                 config.IndexDb,
		MaxRetries:         3,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512 * time.Millisecond,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         0,
		IdleTimeout:        5 * time.Minute,
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
			if _, exists := eventManageFactory.Get(config.EventDestroyPrefix); !exists {
				eventManageFactory.Set(config.EventDestroyPrefix, func(args ...interface{}) {
					_ = client.Close()
					log.Printf("Destroying Redis connection")
				})
			}

			return client
		}

		log.Printf("Failed to connect to Redis, retrying... (attempt %d)", i+1)
		time.Sleep(time.Duration(config.ConnFailRetryInterval) * time.Second)
	}

	panic(fmt.Sprintf("Failed to connect to Redis, reached maximum retry attempts"))
}

// isValidConnection 检查 Redis 连接是否有效
func isValidConnection(client *redis.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	return err == nil
}

// NewClient 初始化 Redis 客户端，并支持多个 Redis 连接
func NewClient(name string) *redis.Client {
	client, loaded := clients.LoadOrStore(name, createClient(name))
	if loaded {
		// 检查连接是否有效
		if !isValidConnection(client.(*redis.Client)) {
			log.Printf("Redis 连接丢失，正在重新连接: %s", name)
			client = createClient(name)
			clients.Store(name, client)
		}
	}
	return client.(*redis.Client)
}
