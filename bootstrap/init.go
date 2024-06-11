package bootstrap

import (
	"os"
	"tool/app/global/variable"
	"tool/app/utils/ants"
	"tool/app/utils/memcached"
	"tool/app/utils/mongo"
	"tool/app/utils/mysql"
	"tool/app/utils/redis"
	"tool/app/utils/yml_config"
	"tool/app/utils/zap_log"

	"go.uber.org/zap"
)

// 初始化加载配置
func Initialize() {

	// 加载配置
	configName := "config"

	//根据环境变量加载不同的配置文件
	if os.Getenv("APP_DEBUG") == "false" {
		configName = "config_production"
	}

	variable.ConfigYml = yml_config.InitViper(configName)

	checkDir()

	//加载日志
	variable.Logs = zap_log.ZapInit(zap_log.ZapLogHandler)
	
	// 创建一个 Ants 池
	pool, _ := ants.NewAnts(variable.ConfigYml.GetInt("Ants.PoolSize"))

	variable.Pool = pool
}

// 检查目录是否存在
func checkDir() {
	// 判断日志目录是否存在 logs
	if _, err := os.Stat(variable.BasePath + "/logs"); os.IsNotExist(err) {
		_ = os.Mkdir(variable.BasePath+"/logs", os.ModePerm)
	}

	// 判断配置文件目录是否存在 public
	if _, err := os.Stat(variable.BasePath + "/public"); os.IsNotExist(err) {
		_ = os.Mkdir(variable.BasePath+"/public", os.ModePerm)
	}

}

// 初始化数据库配置
func InitializeDbConfig() {
	//加载redis
	redisConfig := redis.RedisConfig{
		Host:                  variable.ConfigYml.GetString("Redis.Host"),
		Auth:                  variable.ConfigYml.GetString("Redis.Auth"),
		IndexDb:               variable.ConfigYml.GetInt("Redis.IndexDb"),
		PoolSize:              variable.ConfigYml.GetInt("Redis.PoolSize"),
		MinIdleConns:          variable.ConfigYml.GetInt("Redis.MinIdleConns"),
		ConnFailRetryTimes:    variable.ConfigYml.GetInt("Redis.ConnFailRetryTimes"),
		ConnFailRetryInterval: variable.ConfigYml.GetInt("Redis.ConnFailRetryInterval"),
		EventDestroyPrefix:    variable.EventDestroyPrefix,
	}
	variable.Redis = redis.NewClient("Local", redisConfig)

	//加载memcached
	memcachedConfig := memcached.MemcachedConfig{
		Host:               variable.ConfigYml.GetString("Memcached.Host"),
		EventDestroyPrefix: variable.EventDestroyPrefix,
	}

	variable.Memcached = memcached.NewClient("Local", memcachedConfig)

	// 数据库配置
	config1 := mysql.DatabaseConfig{
		User:               variable.ConfigYml.GetString("Mysql.User"),
		Pass:               variable.ConfigYml.GetString("Mysql.Pass"),
		Host:               variable.ConfigYml.GetString("Mysql.Host"),
		Port:               variable.ConfigYml.GetString("Mysql.Port"),
		Database:           variable.ConfigYml.GetString("Mysql.Database"),
		Charset:            variable.ConfigYml.GetString("Mysql.Charset"),
		SetMaxIdleConns:    variable.ConfigYml.GetInt("Mysql.SetMaxIdleConns"),
		SetMaxOpenConns:    variable.ConfigYml.GetInt("Mysql.SetMaxOpenConns"),
		SetConnMaxLifetime: variable.ConfigYml.GetInt("Mysql.SetConnMaxLifetime"),
		EventDestroyPrefix: variable.EventDestroyPrefix,
	}

	//加载orm
	db, err := mysql.NewDBClient("heal", config1)
	if err != nil {
		variable.Logs.Error("Failed to connect to database", zap.Error(err))
	} else {
		variable.Mysql = db
	}

	if variable.ConfigYml.GetBool("MongoDB.Open") {
		// 数据库配置
		mongoConfig := mongo.DatabaseConfig{
			URI:                variable.ConfigYml.GetString("MongoDB.Uri"),
			Database:           variable.ConfigYml.GetString("MongoDB.Database"),
			MaxPoolSize:        uint64(variable.ConfigYml.GetInt("MongoDB.MaxPoolSize")),
			MinPoolSize:        uint64(variable.ConfigYml.GetInt("MongoDB.MinPoolSize")),
			EventDestroyPrefix: variable.EventDestroyPrefix,
		}

		// 初始化mongo
		variable.MongoDB, err = mongo.InitMongo(mongoConfig)
		if err != nil {
			variable.Logs.Error("Failed to connect to MongoDB", zap.Error(err))
		}
	}
}
