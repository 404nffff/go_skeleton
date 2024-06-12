package admin

import (
	"tool/app/http/controller"
	"tool/app/http/controller/admin"
	"tool/app/http/middleware"
)

// RegisterUserRouter 注册用户路由
func RegisterUserRouter() {

	// 用户路由组
	userGroup := Api.Group("/user")

	// 用户登录
	userGroup.GET("/login", admin.Test)

	userGroup.Use(middleware.AuthMiddleware()) // 将中间件应用于这个路由组
	{
		// 用户登录
		userGroup.GET("/index", controller.Test)
	}

}
