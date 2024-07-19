package admin

import (
	"net/http"

	"tool/global/utils/common"
	"tool/global/variable"
	"tool/pkg/session"
	"tool/server/http/service/admin"

	"github.com/gin-gonic/gin"
)

// Login 登录页面
func LoginPage(c *gin.Context) {

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录",
	})
}

// Login 登录页面提交
func LoginSubmit(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	params := map[string]any{
		"username": username,
		"password": password,
	}

	result, _ := variable.Pool.SubmitTask(admin.Login, params)

	if result["code"] != 200 {
		common.Fail(c, http.StatusBadRequest, result["msg"].(string), nil)
		return
	}

	session.SetM(c, "user", result["data"])

	common.Success(c, "登录成功", nil)
}
