package yml_config

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"tool/global/variable"
	"tool/pkg/yml_config/ymlconfig_interf"

	"github.com/spf13/viper"
)

type ymlConfig struct {
	viper *viper.Viper
}

var configLock sync.Mutex

// 加载配置文件
func LoadConfig(configName string) ymlconfig_interf.YmlConfigInterf {

	//加锁
	configLock.Lock()

	defer configLock.Unlock()

	basePath := variable.BasePath

	v := viper.New()

	// 设置配置文件名称（不带扩展名）
	v.SetConfigName(configName)
	// 设置配置文件类型
	v.SetConfigType("yaml")
	// 设置配置文件路径，可以设置多个路径
	v.AddConfigPath(basePath + "/config")

	// 读取配置文件并处理错误
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件出错: %s", err))
	}

	return &ymlConfig{viper: v}
}

// 封装 viper.GetInt 方法，如果键不存在，则返回默认值
func (y *ymlConfig) GetInt(key string) int {
	if !y.viper.IsSet(key) {
		return 0
	}
	return y.viper.GetInt(key)
}

// 封装 viper.GetString 方法，如果键不存在，则返回默认值
func (y *ymlConfig) GetString(key string) string {
	if !y.viper.IsSet(key) {
		return ""
	}
	return y.viper.GetString(key)
}

// 封装 viper.GetBool 方法，如果键不存在，则返回默认值
func (y *ymlConfig) GetBool(key string) bool {
	if !y.viper.IsSet(key) {
		return false
	}
	return y.viper.GetBool(key)
}

// 封装 viper.GetFloat64 方法，如果键不存在，则返回默认值
func (y *ymlConfig) GetFloat64(key string) float64 {
	if !y.viper.IsSet(key) {
		return 0.0
	}
	return y.viper.GetFloat64(key)
}

// GetDuration 时间单位格式返回值
func (y *ymlConfig) GetDuration(keyName string) time.Duration {

	value := y.viper.GetDuration(keyName)

	return value
}

// GetConfig 根据类型获取配置值，如果键不存在，返回对应类型的默认值
func (y *ymlConfig) GetConfig(key string, defaultValue interface{}) interface{} {
	if !y.viper.IsSet(key) {
		return defaultValue
	}

	switch reflect.TypeOf(defaultValue).Kind() {
	case reflect.Int:
		return y.viper.GetInt(key)
	case reflect.String:
		return y.viper.GetString(key)
	case reflect.Bool:
		return y.viper.GetBool(key)
	case reflect.Float64:
		return y.viper.GetFloat64(key)
	default:
		return defaultValue
	}
}
