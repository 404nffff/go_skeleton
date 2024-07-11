package db_client

import (
	pkgRedis "tool/pkg/redis"

	"github.com/go-redis/redis/v8"
)

func RedisLocal() *redis.Client {

	return pkgRedis.NewClient("Local")
}
