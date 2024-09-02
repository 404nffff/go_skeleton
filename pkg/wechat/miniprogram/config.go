package miniprogram

import "tool/pkg/yml_config"

// LogConfig 日志配置
type LogConfig struct {
	Level string
	Path  string
}

//RedisConfig redis配置
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

	if config.GetString(name+".AppID") == "" {
		panic("Failed to get wechat config")
	}

	return &Config{
		AppID:  config.GetString(name + ".AppID"),
		Secret: config.GetString(name + ".Secret"),
		Debug:  config.GetBool(name + ".Debug"),
		Redis: RedisConfig{
			Addr:     config.GetString(name + ".Redis.Addr"),
			Password: config.GetString(name + ".Redis.Password"),
			DB:       config.GetInt(name + ".Redis.DB"),
		},
		Log: LogConfig{
			Level: config.GetString(name + ".Log.Level"),
			Path: config.GetString(name + ".Log.Path"),
		},
	}
}
