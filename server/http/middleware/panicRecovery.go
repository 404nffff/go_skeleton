package middleware

import (
	"fmt"
	"net/http"
	"tool/global/variable"

	"github.com/gin-gonic/gin"
)

// PanicRecoveryMiddleware 中间件函数
func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误日志
				variable.Logs.Error(fmt.Sprintf("panic info: %v", err))

				// 返回 500 错误给客户端
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Internal server error: %s", err),
				})

				// 中止请求
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}
