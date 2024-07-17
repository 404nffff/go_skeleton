package api

import (
	"fmt"
	"tool/app/global/variable"
	"tool/app/http/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var api *gin.Engine

func InitRouter() *gin.Engine {

	//判断是否是调试模式

	AppDebug := variable.ConfigYml.GetBool("AppDebug")

	if AppDebug == true {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		api = gin.New()
		pprof.Register(api)
	} else {
		gin.SetMode(gin.ReleaseMode)
		api = gin.New()
	}

	// 初始化中间件
	initMiddleware()

	// 初始化路由
	initRoutes()

	return api
}

// 初始化中间件
func initMiddleware() {

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		api.Use(middleware.Cors())
	}

	//使用 gin.Recovery() 中间件
	api.Use(gin.Recovery())

	//使用 LoggerMiddleware 中间件
	logger, err := middleware.InitLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	defer logger.Sync()

	//初始化日志
	api.Use(middleware.LoggerMiddleware(logger))

	//初始化session
	api.Use(middleware.SessionMiddleware())

	//初始化参数验证
	api.Use(middleware.ValidateParams())
}
