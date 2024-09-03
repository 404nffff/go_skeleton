package api

import (
	"tool/pkg/web_server"
	"tool/server/http/controller/tool"

	"github.com/gin-gonic/gin"
)

// 注册路由
func init() {

	web_server.RegisterRoutes("/tool/oss",
		web_server.Route{
			Method:   "POST",
			Path:     "/upload",
			Handlers: []gin.HandlerFunc{tool.Upload},
		},
	)
}
