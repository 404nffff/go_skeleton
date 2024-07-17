package admin

import (
	"tool/app/http/controller/admin"
	"tool/app/http/middleware"
)

// RegisterUserRouter 注册用户路由
func RegisterUserRouter() {

	adminGroup := Api.Group("/admin")

	adminGroup.Use(middleware.AdminAuthMiddleware()) // 将中间件应用于这个路由组
	{
		//后台首页
		adminGroup.GET("/index", admin.Index)
	}

}
