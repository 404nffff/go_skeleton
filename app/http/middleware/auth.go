package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware : 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization
		tokenString := c.GetHeader("Authorization")
		// 验证 token 格式
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请求头中auth为空"})
			c.Abort()
			return
		}

		// token 校验
		if tokenString != "Bearer 123456" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}

		c.Next()
	}
}
