package memcached

import (
	"fmt"
	"log"
	"sync"
	"tool/pkg/event_manage"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	clients sync.Map // 用于存储多个 Memcached 客户端实例
	mu      sync.Mutex
)

type MemcachedConfig struct {
	Host               string // Memcached 服务器地址，格式为 "host:port"。
	EventDestroyPrefix string // 事件销毁前缀
}

// NewClient 创建一个新的 Memcached 客户端
// 使用 sync.Once 确保只初始化一次
func NewClient(name string, config MemcachedConfig) *memcache.Client {
	client, _ := clients.LoadOrStore(name, createClient(name, config))
	return client.(*memcache.Client)
}

// createClient 创建一个新的 Memcached 客户端
func createClient(name string, config MemcachedConfig) *memcache.Client {
	mu.Lock()
	defer mu.Unlock()

	var client *memcache.Client
	var err error

	client = memcache.New(config.Host)

	err = Ping(client)
	if err != nil {
		log.Fatalf("Failed to connect to Memcached: %v", err)
	}

	// 注册销毁事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	if _, exists := eventManageFactory.Get(config.EventDestroyPrefix + "Memcached_" + name); exists == false {
		eventManageFactory.Set(config.EventDestroyPrefix+"Memcached_"+name, func(args ...interface{}) {
			// Memcached 客户端不需要显式关闭连接
			log.Printf("Destroying Memcached connection for %s", name)
			client.Close()
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
