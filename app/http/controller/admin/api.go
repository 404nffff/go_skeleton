package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {

	fmt.Println("admin test")

	c.HTML(http.StatusOK, "base", gin.H{
		"title":   "User Page",
		"content": "user.tmpl", // 指定内容模板
	})
}
