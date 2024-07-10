package api

import (
	"tool/app/http/controller"
	"tool/app/http/middleware"
)

// RegisterUserRouter 注册用户路由
func RegisterUserRouter() {

	// 用户路由组
	userGroup := Api.Group("/user")

	// 用户登录
	userGroup.GET("/test", controller.Test)
	userGroup.GET("/test2", controller.Test2)
	userGroup.GET("/user", controller.GetUsersHandler)
	userGroup.GET("/user2", controller.GetUserMongo)
	userGroup.GET("/testudp", controller.TestUdp)
	userGroup.GET("/task", controller.TestAnt)
	userGroup.GET("/s", controller.StatusAnt)
	userGroup.POST("/upload", controller.Upload)
	
	userGroup.Use(middleware.AuthMiddleware()) // 将中间件应用于这个路由组
	{
		// 用户登录
		userGroup.GET("/login", controller.Test)
	}

}
