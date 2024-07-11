package bootstrap

import (
	"os"
	"tool/app/global/variable"
	"tool/pkg/ants"
	"tool/pkg/yml_config"
	"tool/pkg/zap_log"
)

// 初始化加载配置
func Initialize() {

	// 加载配置
	configName := "config"

	//根据环境变量加载不同的配置文件
	if os.Getenv("APP_DEBUG") == "false" {
		configName = "config_production"
	}

	variable.ConfigYml = yml_config.LoadConfig(configName)

	checkDir()

	//加载日志
	variable.Logs = zap_log.ZapInit(zap_log.ZapLogHandler)
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

// 初始化协程池
func InitPool(poolSize int) {
	// 创建一个 Ants 池
	pool, _ := ants.NewAnts(poolSize)

	variable.Pool = pool
}
