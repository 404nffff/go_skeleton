package admin

import (
	"tool/app/http/controller/admin"
)

// RegisterIndexRouter 注册用户路由
func RegisterRouter() {

	// 用户登录
	Api.GET("/login", admin.Test)

	RegisterUserRouter()
}
