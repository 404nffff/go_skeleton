package miniprogram

import (
	"tool/pkg/yml_config"
)

// LogConfig 日志配置
type LogConfig struct {
	Level string
	Path  string
}

// RedisConfig redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Config 配置
type Config struct {
	AppID  string
	Secret string
	Debug  bool
	Redis  RedisConfig
	Log    LogConfig
}

func loadConfig(name string) *Config {
	config := yml_config.LoadConfig("wechat")

	if config.GetString("MiniPro."+name+".AppID") == "" {
		panic("Failed to get wechat config")
	}

	return &Config{
		AppID:  config.GetString("MiniPro." + name + ".AppID"),
		Secret: config.GetString("MiniPro." + name + ".AppSecret"),
		Debug:  config.GetBool("MiniPro." + name + ".Debug"),
		Redis: RedisConfig{
			Addr:     config.GetString("MiniPro." + name + ".Redis.Addr"),
			Password: config.GetString("MiniPro." + name + ".Redis.Password"),
			DB:       config.GetInt("MiniPro." + name + ".Redis.DB"),
		},
		Log: LogConfig{
			Level: config.GetString("MiniPro." + name + ".Log.Level"),
			Path:  config.GetString("MiniPro." + name + ".Log.Path"),
		},
	}
}
