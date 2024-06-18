package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tool/app/global/variable"
	"tool/app/utils/event_manage"
	"tool/app/utils/process"
	"tool/app/utils/tcp"
	"tool/bootstrap"
	"tool/routers/admin"

	"go.uber.org/zap"
)

var (
	server *http.Server
)

func init() {
	bootstrap.Initialize()

	//初始化协程池
	bootstrap.InitPool(variable.ConfigYml.GetInt("HttpServer.Admin.WorkNum"))
}

func main() {
	process.Initialize("admin", startServerInForeground)
}

func initWebServer() *http.Server {
	port := variable.ConfigYml.GetString("HttpServer.Admin.Port")

	if tcp.IsPortInUse(port) {
		variable.Logs.Fatal("Port is already in use:", zap.String("port", port))
	}

	variable.Logs.Info("Starting server on port", zap.String("port", port))

	router := admin.InitRouter()

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

	bootstrap.InitializeDbConfig()

	go destroy()

	defer variable.Pool.Release()

	server = initWebServer()
	variable.Logs.Info("Starting server in foreground...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		variable.Logs.Fatal("Listen error:", zap.Error(err))
	}
}

func destroy() {

	//打印日志
	variable.Logs.Info("Destroying server...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		received := <-c
		variable.Logs.Warn("ProcessKilled", zap.String("信号值", received.String()))
		(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)
		os.Exit(1)
	}()
}
