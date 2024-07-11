package db_client

import (
	"tool/pkg/memcached"

	"github.com/bradfitz/gomemcache/memcache"
)

// 本地连接 Memcached
func MemLocal() *memcache.Client {
	return memcached.NewClient("Local")
}
