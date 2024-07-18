package api

import (
	"tool/app/http/controller/tool"
	"tool/pkg/web_server"
)

// 注册路由
func init() {

	web_server.RegisterRoutes("/tool/oss",
		web_server.Route{
			Method:  "POST",
			Path:    "/upload",
			Handler: tool.Upload,
		},
	)
}
