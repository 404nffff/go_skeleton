package memcached

import (
	"fmt"
	"log"
	"sync"
	"time"
	"tool/pkg/event_manage"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	clients sync.Map // 用于存储多个 Memcached 客户端实例
)

// NewClient 创建一个新的 Memcached 客户端
func NewClient(name string) *memcache.Client {
	client, loaded := clients.LoadOrStore(name, createClient(name))

	if loaded {
		// 检查连接是否有效
		if !isValidConnection(client.(*memcache.Client)) {
			log.Printf("Memcached 连接丢失，正在重新连接: %s", name)
			client = createClient(name)
			clients.Store(name, client)
		}
	}

	return client.(*memcache.Client)
}

// createClient 创建一个新的 Memcached 客户端
func createClient(name string) *memcache.Client {
	// 加载配置
	config := loadConfig(name)

	// 判断是否为空
	if config.Host == "" {
		panic(fmt.Sprintf("Failed to get Memcached config: %s", name))
	}

	maxRetries := config.ConnFailRetryTimes                                    // 最大重试次数
	retryInterval := time.Duration(config.ConnFailRetryInterval) * time.Second // 重试间隔

	var client *memcache.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client = memcache.New(config.Host)
		err = Ping(client)
		if err == nil {
			break // 连接成功，跳出循环
		}
		log.Printf("Failed to connect to Memcached, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Memcached after %d attempts: %v", maxRetries, err))
	}

	// 注册销毁事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	if _, exists := eventManageFactory.Get(config.EventDestroyPrefix); !exists {
		eventManageFactory.Set(config.EventDestroyPrefix, func(args ...interface{}) {
			log.Printf("Destroying Memcached connection for %s", name)
			client = nil
		})
	}

	return client
}

// Ping 检测 Memcached 连接是否连通
func Ping(client *memcache.Client) error {
	key := "ping_key"
	value := []byte("ping_value")

	// 尝试设置一个键值对
	err := client.Set(&memcache.Item{Key: key, Value: value, Expiration: 10})
	if err != nil {
		return err
	}

	// 尝试获取刚刚设置的键值对
	_, err = client.Get(key)
	if err != nil {
		return err
	}

	return nil
}

// isValidConnection 检查 Memcached 连接是否有效
func isValidConnection(client *memcache.Client) bool {
	err := Ping(client)
	return err == nil
}

// Set 设置一个键值对
func Set(clientName, key string, value []byte, expiration int32) error {
	client, err := getClient(clientName)
	if err != nil {
		return err
	}

	err = client.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
	}
	return err
}

// Get 获取一个键值对
func Get(clientName, key string) ([]byte, error) {
	client, err := getClient(clientName)
	if err != nil {
		return nil, err
	}

	item, err := client.Get(key)
	if err != nil {
		log.Printf("Failed to get key %s: %v", key, err)
		return nil, err
	}
	return item.Value, nil
}

// Delete 删除一个键值对
func Delete(clientName, key string) error {
	client, err := getClient(clientName)
	if err != nil {
		return err
	}

	err = client.Delete(key)
	if err != nil {
		log.Printf("Failed to delete key %s: %v", key, err)
	}
	return err
}

// getClient 获取指定名称的 Memcached 客户端
func getClient(name string) (*memcache.Client, error) {
	client, ok := clients.Load(name)
	if !ok {
		return nil, ErrClientNotInitialized
	}
	return client.(*memcache.Client), nil
}

// ErrClientNotInitialized is returned when the Memcached client is not initialized
var ErrClientNotInitialized = fmt.Errorf("Memcached client not initialized")
