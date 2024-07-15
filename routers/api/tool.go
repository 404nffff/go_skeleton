package api

import "tool/app/http/controller/tool"

// 注册路由
func init() {

	registerRoutesToGroup("/tool/oss",
		route{
			method:  "POST",
			path:    "/upload",
			handler: tool.Upload,
		},
	)
}
