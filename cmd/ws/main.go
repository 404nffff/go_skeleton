package main

import (
	"time"
	"tool/bootstrap"
	"tool/global/variable"
	"tool/pkg/event_manage"
	"tool/pkg/process"
	"tool/pkg/web_server"

	"github.com/gin-gonic/gin"

	"tool/server/websocket/handle"
	_ "tool/server/websocket/routers" // 加载api路由
)

func init() {
	bootstrap.Initialize()

	//初始化协程池
	bootstrap.InitPool(variable.ConfigYml.GetInt("HttpServer.Ws.WorkNum"))
}

func main() {
	process.Initialize("ws", startServerInForeground)
}

func startServerInForeground() {

	config := web_server.RouterConfig{
		AppDebug:          variable.ConfigYml.GetBool("AppDebug"),
		CustomMiddlewares: []gin.HandlerFunc{
			//middleware.SessionMiddleware(),
			//middleware.ValidateParams(),
		},
	}

	// 初始化路由
	router := web_server.InitRouter(config)

	port := variable.ConfigYml.GetString("HttpServer.Ws.Port")

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

	go handle.HandleMsg()

	server.Start()

}
