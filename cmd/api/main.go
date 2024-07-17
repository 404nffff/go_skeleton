package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"tool/app/global/variable"
	"tool/app/http/routers/api"
	"tool/bootstrap"
	"tool/pkg/event_manage"
	"tool/pkg/process"
	"tool/pkg/tcp"

	"go.uber.org/zap"
)

var (
	server *http.Server
)

func init() {
	bootstrap.Initialize()

	//初始化协程池
	bootstrap.InitPool(variable.ConfigYml.GetInt("HttpServer.Api.WorkNum"))
}

func main() {
	process.Initialize("api", startServerInForeground)
}

func initWebServer() *http.Server {
	port := variable.ConfigYml.GetString("HttpServer.Api.Port")

	if tcp.IsPortInUse(port) {
		variable.Logs.Fatal("Port is already in use:", zap.String("port", port))
	}

	variable.Logs.Info("Starting server on port", zap.String("port", port))

	router := api.InitRouter()

	server = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server
}

func startServerInForeground() {

	//开启tcp协议
	//go tcp.NewTCPServer(":8081").Start()

	go destroy()

	defer variable.Pool.Release()

	server = initWebServer()
	variable.Logs.Info("Starting server in foreground...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		variable.Logs.Fatal("Listen error:", zap.Error(err))
	}
}

func destroy() {
	// 打印日志
	variable.Logs.Info("Destroying server...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		received := <-c
		variable.Logs.Warn("ProcessKilled", zap.String("信号值", received.String()))
		(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)

		// 等待所有任务完成
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			variable.Pool.Release()
		}()
		wg.Wait()

		// 关闭服务器
		if err := server.Shutdown(context.Background()); err != nil {
			variable.Logs.Fatal("Server forced to shutdown:", zap.Error(err))
		}

		os.Exit(1)
	}()
}
