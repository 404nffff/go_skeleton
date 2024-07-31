package routers

import (
	"tool/pkg/web_server"
	"tool/server/websocket/handle"
)

// 注册路由
func init() {

	web_server.RegisterRoutes("",
		web_server.Route{
			Method:  "GET",
			Path:    "/join",
			Handler: handle.Join,
		},
	)
}
