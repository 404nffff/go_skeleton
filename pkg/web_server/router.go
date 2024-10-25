package web_server

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Route 表示单个路由的结构体
type Route struct {
	Method      string            // HTTP 方法（GET, POST 等）
	Path        string            // 路由路径
	Handlers    []gin.HandlerFunc // 处理函数列表
	Middlewares []gin.HandlerFunc // 路由特定的中间件
	Params      reflect.Type      // 路由参数
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
	groups map[string]*RouterGroup
}

// RouterGroup 表示一个路由组
type RouterGroup struct {
	prefix        string
	middlewares   []gin.HandlerFunc
	routes        []Route
	beforeHandler []Route // 前置处理函数
	afterHandler  []Route // 后置处理函数
}

var (
	globalRouter     *Router                      // 全局路由器实例
	once             sync.Once                    // 确保只初始化一次
	routeGroup       map[string][]Route           // 存储路由组
	middlewaresGroup map[string][]gin.HandlerFunc // 存储中间件组
	RouteParamMap    map[string]reflect.Type      // 路径到结构体类型的映射
)

func init() {
	once.Do(func() {
		routeGroup = make(map[string][]Route)
		middlewaresGroup = make(map[string][]gin.HandlerFunc)
		RouteParamMap = make(map[string]reflect.Type)
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

			if route.Params != nil {

				// 将路由路径和参数类型映射起来
				RouteParamMap[prefix+route.Path] = route.Params
			}

			handlers := append(route.Middlewares, route.Handlers...)
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

// Group 在现有路由组中创建一个子组
func Group(prefix string, middlewares ...gin.HandlerFunc) RouterGroup {
	newGroup := RouterGroup{
		prefix:      prefix,
		middlewares: middlewares,
	}

	return newGroup
}

// Add 添加一个路由到路由组
func (rg *RouterGroup) Add(method, path string, handlers ...gin.HandlerFunc) *Route {
	route := Route{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	}
	rg.routes = append(rg.routes, route)
	return &rg.routes[len(rg.routes)-1]
}

// GET 是添加 GET 路由的快捷方法
func (rg *RouterGroup) GET(path string, handlers ...gin.HandlerFunc) *Route {
	return rg.Add("GET", path, handlers...)
}

// POST 是添加 POST 路由的快捷方法
func (rg *RouterGroup) POST(path string, handlers ...gin.HandlerFunc) *Route {
	return rg.Add("POST", path, handlers...)
}

// PUT 是添加 PUT 路由的快捷方法
func (rg *RouterGroup) PUT(path string, handlers ...gin.HandlerFunc) *Route {
	return rg.Add("PUT", path, handlers...)
}

// DELETE 是添加 DELETE 路由的快捷方法
func (rg *RouterGroup) DELETE(path string, handlers ...gin.HandlerFunc) *Route {
	return rg.Add("DELETE", path, handlers...)
}

// 添加前置处理函数
func (rg *RouterGroup) Before(handlers ...gin.HandlerFunc) *RouterGroup {

	route := Route{
		Handlers: append([]gin.HandlerFunc{}, handlers...),
	}

	rg.beforeHandler = append(rg.beforeHandler, route)
	return rg
}

// 添加后置处理函数
func (rg *RouterGroup) After(handlers ...gin.HandlerFunc) *RouterGroup {

	route := Route{
		Handlers: append([]gin.HandlerFunc{}, handlers...),
	}

	rg.afterHandler = append(rg.afterHandler, route)
	return rg
}

// Params 设置路由参数类型
func (r *Route) Bind(params interface{}) *Route {
	r.Params = reflect.TypeOf(params)
	return r
}

// Use 添加中间件到路由
func (r *Route) Use(middlewares ...gin.HandlerFunc) *Route {
	r.Middlewares = append(r.Middlewares, middlewares...)
	return r
}

// ImportRoutes 导入路由到全局路由器
func ImportRoutes(routerGroups ...RouterGroup) {
	// 为每个路由组添加路由
	for _, rg := range routerGroups {

		if rg.middlewares != nil {
			//注册中间件
			middlewaresGroup[rg.prefix] = append(middlewaresGroup[rg.prefix], rg.middlewares...)
		}

		

		for _, route := range rg.routes {

			
			//注册前置处理函数、把handlers添加到handlers前面
			if len(rg.beforeHandler) > 0 {
				for _, before := range rg.beforeHandler {
					route.Handlers = append(before.Handlers, route.Handlers...)
				}
			}

			//注册后置处理函数、把handlers添加到handlers后面
			if len(rg.afterHandler) > 0 {
				for _, after := range rg.afterHandler {
					route.Handlers = append(route.Handlers, after.Handlers...)
				}
			}

			//注册路由
			routeGroup[rg.prefix] = append(routeGroup[rg.prefix], route)
		}

		
	}
}
