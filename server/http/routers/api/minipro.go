package api

import (
	"tool/pkg/web_server"
	"tool/server/http/controller/wechat"

	"github.com/gin-gonic/gin"
)

// 注册路由 - 小程序
func init() {

	web_server.RegisterRoutes("/front/minipro",
		web_server.Route{
			Method:   "POST",
			Path:     "/auth",
			Handlers: []gin.HandlerFunc{wechat.Auth},
			//Params:   reflect.TypeOf(request.AuthParams{}),
		},
	)
}
