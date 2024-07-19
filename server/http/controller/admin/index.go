package admin

import (
	"fmt"
	"net/http"
	"tool/pkg/session"

	"github.com/gin-gonic/gin"
)

// func Index(c *gin.Context) {

// 	sessionData := session.GetM(c, "user", "username")

// 	fmt.Println(sessionData)

//		c.HTML(http.StatusOK, "base", gin.H{
//			"title":   "管理系统",
//			"content": "admin/index.html", // 指定内容模板
//		})
//	}
func Index(c *gin.Context) {

	sessionData := session.GetM(c, "user", "username")

	fmt.Println(sessionData)

	c.HTML(http.StatusOK, "adminv2/index.html", gin.H{
		"title": "管理系统",
	})
}
