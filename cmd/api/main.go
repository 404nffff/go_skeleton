package main

import (
	"time"
	"tool/bootstrap"
	"tool/global/variable"
	"tool/pkg/event_manage"
	"tool/pkg/process"
	"tool/pkg/web_server"
	"tool/server/http/middleware"

	"github.com/gin-gonic/gin"

	_ "tool/server/http/routers/api" // 加载api路由
)

func init() {
	bootstrap.Initialize()

	//初始化协程池
	bootstrap.InitPool(variable.ConfigYml.GetInt("HttpServer.Api.WorkNum"))
}

func main() {
	process.Initialize("api", startServerInForeground)
}

func startServerInForeground() {

	config := web_server.RouterConfig{
		AppDebug:         variable.ConfigYml.GetBool("AppDebug"),
		AllowCrossDomain: variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain"),
		CustomMiddlewares: []gin.HandlerFunc{
			middleware.SessionMiddleware(),
			middleware.ValidateParams(),
			middleware.Cors(),
		},
	}

	// 初始化路由
	router := web_server.InitRouter(config)

	port := variable.ConfigYml.GetString("HttpServer.Api.Port")

	webConfig := web_server.ServerConfig{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Logger:         variable.Logs,
		DestroyCallback: func() {

			variable.Pool.Release()

			//睡 1 秒，等待所有协程执行完毕
			time.Sleep(1 * time.Second)

			// 自定义的销毁逻辑
			(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)
		},
	}

	server := web_server.NewServer(webConfig)

	server.Start()

}
