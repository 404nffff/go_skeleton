package admin

import (
	"tool/server/http/controller/admin"
	"tool/server/http/middleware"
)

// RegisterIndexRouter 注册用户路由
func RegisterRouter() {

	// 用户登录
	Api.GET("/admin/login", middleware.RateLimiter(10, 1), admin.LoginPage)

	// 用户登录提交
	Api.POST("/admin/login_submit", middleware.RateLimiter(10, 1), admin.LoginSubmit)

	RegisterUserRouter()
}
