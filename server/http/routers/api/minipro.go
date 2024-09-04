package api

import (
	"reflect"
	"tool/pkg/web_server"
	"tool/server/http/controller/wechat"
	"tool/server/http/request/wechat/minipro"

	"github.com/gin-gonic/gin"
)

// 注册路由 - 小程序
func init() {

	web_server.RegisterRoutes("/front/minipro",
		web_server.Route{
			Method:   "POST",
			Path:     "/auth",
			Handlers: []gin.HandlerFunc{wechat.Auth},
			Params:   reflect.TypeOf(minipro.AuthParams{}),
		},
	)
}
