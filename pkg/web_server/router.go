package web_server

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Route 表示单个路由的结构体
type Route struct {
	Method      string            // HTTP 方法（GET, POST 等）
	Path        string            // 路由路径
	Handler     gin.HandlerFunc   // 处理函数
	Middlewares []gin.HandlerFunc // 路由特定的中间件
}

// RouterConfig 保存路由器的配置
type RouterConfig struct {
	AppDebug          bool              // 是否开启调试模式
	CustomMiddlewares []gin.HandlerFunc // 自定义中间件列表
}

// Router 封装 gin.Engine 和相关配置
type Router struct {
	engine *gin.Engine  // gin 引擎实例
	config RouterConfig // 路由器配置
	logger *zap.Logger  // zap 日志实例
}

var (
	globalRouter     *Router                      // 全局路由器实例
	once             sync.Once                    // 确保只初始化一次
	routeGroup       map[string][]Route           // 存储路由组
	middlewaresGroup map[string][]gin.HandlerFunc // 存储中间件组
)

func init() {
	once.Do(func() {
		routeGroup = make(map[string][]Route)
		middlewaresGroup = make(map[string][]gin.HandlerFunc)
	})
}

// GetGlobalRouter 获取全局路由器实例
func GetGlobalRouter() *Router {
	if globalRouter == nil {
		panic("Global router not initialized. Call InitGlobalRouter first.")
	}
	return globalRouter
}

// Init 初始化路由
func (r *Router) Init() *gin.Engine {
	if r.config.AppDebug {
		r.engine = gin.New()
		pprof.Register(r.engine)
	} else {
		gin.SetMode(gin.ReleaseMode)
		r.engine = gin.New()
	}

	r.initMiddleware()
	r.initRoutes()

	return r.engine
}

// initMiddleware 初始化中间件
func (r *Router) initMiddleware() {
	// 初始化日志
	logger, err := InitLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	r.logger = logger

	r.engine.Use(gin.Recovery())
	r.engine.Use(LoggerMiddleware(r.logger))

	// 添加自定义中间件
	for _, m := range r.config.CustomMiddlewares {
		r.engine.Use(m)
	}
}

// initRoutes 初始化路由
func (r *Router) initRoutes() {
	for prefix, routes := range routeGroup {
		var group *gin.RouterGroup

		// 检查是否有该前缀的中间件
		if middlewares, exists := middlewaresGroup[prefix]; exists && len(middlewares) > 0 {
			// 如果存在中间件，使用它们创建组
			group = r.engine.Group(prefix, middlewares...)
		} else {
			// 如果不存在中间件，创建没有中间件的组
			group = r.engine.Group(prefix)
		}

		// 为组添加路由
		for _, route := range routes {
			handlers := []gin.HandlerFunc{route.Handler}

			// 检查路由是否有中间件
			if len(route.Middlewares) > 0 {
				handlers = append(route.Middlewares, route.Handler)
			}

			group.Handle(route.Method, route.Path, handlers...)
		}
	}
}

// GetLogger 获取日志器
func (r *Router) GetLogger() *zap.Logger {
	return r.logger
}

// 以下是公共方法，用于操作全局路由器实例

// RegisterRoutes 注册路由到全局路由器
func RegisterRoutes(prefix string, routes ...Route) {
	for _, route := range routes {
		routeGroup[prefix] = append(routeGroup[prefix], route)
	}
}

// RegisterMiddleware 注册中间件到全局路由器
func RegisterMiddleware(prefix string, middlewares ...gin.HandlerFunc) {

	for _, middleware := range middlewares {
		middlewaresGroup[prefix] = append(middlewaresGroup[prefix], middleware)
	}
}

// InitRouter 初始化并返回全局路由器的 gin.Engine
func InitRouter(config RouterConfig) *gin.Engine {

	globalRouter = &Router{
		config: config,
	}

	return GetGlobalRouter().Init()
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() *zap.Logger {
	return GetGlobalRouter().GetLogger()
}
