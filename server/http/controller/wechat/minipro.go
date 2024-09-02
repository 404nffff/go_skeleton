package wechat

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthParams struct {
	Code string `json:"code"  binding:"required"`
}

func Auth(c *gin.Context) {

	//miniprogramClient := miniprogram.NewMiniProgramClient("Default")

	var params AuthParams
	c.ShouldBindJSON(params)

	fmt.Println(params)

	// ctx := c.Request.Context()

	// // 获取用户信息
	// userInfo, err := miniprogramClient.Auth.Session(ctx, code)

	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"code": 1,
	// 		"msg":  err.Error(),
	// 	})
	// 	return
	// }

	// fmt.Println(userInfo)
}
