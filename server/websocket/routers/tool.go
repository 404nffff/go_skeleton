package routers

import (
	"tool/pkg/web_server"
	"tool/server/websocket/handle"

	"github.com/gin-gonic/gin"
)

// 注册路由
func init() {

	web_server.RegisterRoutes("",
		web_server.Route{
			Method:   "GET",
			Path:     "/join",
			Handlers: []gin.HandlerFunc{handle.Join},
		},
	)
}
