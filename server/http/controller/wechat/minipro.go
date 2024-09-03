package wechat

import (
	"fmt"
	"tool/server/http/request"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {

	params, _ := c.Get("params")

	fmt.Println(params.(*request.AuthParams).Code)

	//miniprogramClient := miniprogram.NewMiniProgramClient("Default")

	// var params AuthParams

	// if error := c.ShouldBindJSON(&params); error != nil {
	// 	c.JSON(200, gin.H{
	// 		"code": 1,
	// 		"msg":  error.Error(),
	// 	})
	// 	return
	// }

	// fmt.Println(params)

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
