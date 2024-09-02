package api

import (
	"tool/pkg/web_server"
	"tool/server/http/controller/wechat"
)

// 注册路由 - 小程序
func init() {

	web_server.RegisterRoutes("/front/minipro",
		web_server.Route{
			Method:  "POST",
			Path:    "/auth",
			Handler: wechat.Auth,
		},
	)
}
