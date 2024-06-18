package session

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetSession 设置会话
func Set(c *gin.Context, key string, value interface{}) string {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()

	return session.ID()
}

// GetSession 获取会话
func Get(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

// ClearSession 清除会话
func Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// SetMemcached 设置会话
func SetM(c *gin.Context, key string, value interface{}) string {

	//value 转成json string
	valueJson, _ := json.Marshal(value)

	session := sessions.Default(c)
	session.Set(key, string(valueJson))
	session.Save()

	return session.ID()
}

// GetMemcached 获取会话 memcached
// key: 会话key
// item: 会话key中的item
func GetM(c *gin.Context, key string, item string) interface{} {
	session := sessions.Default(c)
	value := session.Get(key)

	if value == nil {
		return nil
	}

	//value 转成json string
	itemValue, err := jsonparser.GetString([]byte(value.(string)), item)
	if err != nil {
		return nil
	}

	return itemValue
}
