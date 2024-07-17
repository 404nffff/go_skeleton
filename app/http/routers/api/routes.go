package api

import "github.com/gin-gonic/gin"

// route 表示单个路由的结构体
type route struct {
    method      string
    path        string
    handler     gin.HandlerFunc
    middlewares []gin.HandlerFunc
}

// routeGroup 表示路由组的结构体
type routeGroup struct {
    prefix      string
    routes      []route
    middlewares []gin.HandlerFunc
}

var groups []routeGroup

// registerRoutesToGroup 注册多个路由到指定路由组
func registerRoutesToGroup(prefix string, routes ...route) {
    // 查找是否已经存在指定前缀的路由组
    for i, group := range groups {
        if group.prefix == prefix {
            // 如果找到了,将路由添加到该路由组
            groups[i].routes = append(groups[i].routes, routes...)
            return
        }
    }

    // 如果没有找到指定前缀的路由组,创建一个新的路由组并添加路由
    groups = append(groups, routeGroup{
        prefix: prefix,
        routes: routes,
    })
}

// registerMiddlewares 注册中间件到指定路由组
func registerMiddlewares(prefix string, middlewares ...gin.HandlerFunc) {
    for i, group := range groups {
        if group.prefix == prefix {
            groups[i].middlewares = append(groups[i].middlewares, middlewares...)
            return
        }
    }
    
    // 如果没有找到指定前缀的路由组,创建一个新的路由组并添加中间件
    groups = append(groups, routeGroup{
        prefix:      prefix,
        middlewares: middlewares,
    })
}

// initRoutes 初始化并注册所有路由
func initRoutes() {
    for _, group := range groups {
        // 创建对应的路由组
        g := api.Group(group.prefix)
        
        // 应用路由组的中间件
        g.Use(group.middlewares...)
        
        // 注册路由组下的所有路由
        for _, route := range group.routes {
            // 创建一个新的路由组,用于应用单个路由的中间件
            subGroup := g.Group("")
            subGroup.Use(route.middlewares...)
            
            // 注册路由
            subGroup.Handle(route.method, route.path, route.handler)
        }
    }
}
