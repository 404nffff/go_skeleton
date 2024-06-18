package middleware

import (
	"tool/app/utils/session"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware : 后台认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取session
		authSession := session.Get(c, "user")

		if authSession == nil {
			//跳转到登录页面
			c.Redirect(302, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
