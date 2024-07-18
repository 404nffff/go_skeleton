package web_server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tool/pkg/tcp"

	"go.uber.org/zap"
)

// DestroyCallback 定义了销毁时的回调函数类型
type DestroyCallback func()

// ServerConfig 定义了HTTP服务器的配置
type ServerConfig struct {
	Addr            string          // 服务器监听的地址
	Handler         http.Handler    // HTTP请求处理器
	ReadTimeout     time.Duration   // 读取请求的最大时间
	WriteTimeout    time.Duration   // 写入响应的最大时间
	MaxHeaderBytes  int             // 请求头的最大字节数
	Logger          *zap.Logger     // Zap日志记录器
	DestroyCallback DestroyCallback // 服务器销毁时的回调函数
}

// HttpServer 表示一个HTTP服务器
type HttpServer struct {
	http            *http.Server    // 底层的http.Server
	logger          *zap.Logger     // Zap日志记录器
	destroyCallback DestroyCallback // 服务器销毁时的回调函数
}

// NewServer 创建一个新的HttpServer实例
func NewServer(config ServerConfig) *HttpServer {
	port := config.Addr

	// 检查端口是否已被占用
	if tcp.IsPortInUse(port) {
		config.Logger.Fatal("Port is already in use", zap.String("port", port))
	}

	// 创建底层的http.Server
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        config.Handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &HttpServer{
		http:            server,
		logger:          config.Logger,
		destroyCallback: config.DestroyCallback,
	}
}

// Start 启动HTTP服务器
func (server *HttpServer) Start() {

	// 在后台启动 HTTP 服务器
	go server.http.ListenAndServe()

	// 创建一个 channel 用于控制程序退出
	// 这确保了主函数不会在清理操作完成之前退出
	done := make(chan bool)

	go func() {
		// 创建一个 channel 来接收操作系统信号
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞直到接收到信号
		<-c
		// 收到信号后，记录日志
		server.logger.Info("Received shutdown signal")
		// 调用 destroy 方法来清理资源
		server.destroy()
		// 发送信号表示清理完成，允许程序退出
		done <- true
	}()

	// 主函数阻塞在这里，直到接收到 done 信号
	<-done
	// 最后的日志，表示服务器已完全退出
	server.logger.Info("Server exited")
}

func (server *HttpServer) destroy() {
	server.logger.Info("Destroying server...")

	// 如果设置了销毁回调函数，执行它
	if server.destroyCallback != nil {
		server.logger.Info("Executing destroy callback")
		server.destroyCallback()
		server.logger.Info("Destroy callback completed")
	}

	server.logger.Info("Shutting down HTTP server")
	// 创建一个带超时的 context，用于 http.Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 确保在函数退出时调用 cancel，以释放资源
	defer cancel()

	// 优雅地关闭 HTTP 服务器
	if err := server.http.Shutdown(ctx); err != nil {
		// 如果关闭过程中出现错误，记录错误日志
		server.logger.Error("Server shutdown error", zap.Error(err))
	} else {
		// 如果成功关闭，记录成功日志
		server.logger.Info("Server shutdown completed successfully")
	}

	server.logger.Info("Server gracefully stopped")
	// 短暂睡眠，给日志系统一些时间来刷新缓冲区
	// 这有助于确保所有日志都被写入，特别是在使用异步日志库时
	time.Sleep(100 * time.Millisecond)
}
