package ymlconfig_interf

import "time"

// YmlConfigInterf 定义接口
type YmlConfigInterf interface {
	GetInt(key string) int
	GetString(key string) string
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetConfig(key string, defaultValue interface{}) interface{}
	GetDuration(keyName string) time.Duration
}
