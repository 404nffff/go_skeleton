package middleware

import (
	"tool/app/global/variable"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memcached"
	"github.com/gin-gonic/gin"
)

// SessionMiddleware 初始化会话中间件
func SessionMiddleware() gin.HandlerFunc {

	secret := variable.ConfigYml.GetString("Session.Secret")
	maxAge := variable.ConfigYml.GetInt("Session.MaxAge")
	name := variable.ConfigYml.GetString("Session.Name")

	saveMethod := variable.ConfigYml.GetString("Session.SaveMethod")

	var store sessions.Store

	// 保存方式
	if saveMethod == "cookie" {
		store = cookie.NewStore([]byte(secret))

		store.Options(sessions.Options{
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
		})
	} else if saveMethod == "memcached" {
		store = memcached.NewStore(variable.Memcached, "", []byte(secret))
	}

	return sessions.Sessions(name, store)
}
