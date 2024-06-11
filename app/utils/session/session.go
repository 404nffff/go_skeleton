package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetSession 设置会话
func SetSession(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

// GetSession 获取会话
func GetSession(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

// ClearSession 清除会话
func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
