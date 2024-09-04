package wechat

import (
	"fmt"
	"tool/global/utils/common"
	"tool/global/utils/wechat"
	"tool/server/http/request/wechat/minipro"

	"github.com/gin-gonic/gin"
)

// Auth 小程序登录
func Auth(c *gin.Context) {

	params, _ := c.Get("params")

	code := params.(*minipro.AuthParams).Code

	miniprogramClient := wechat.MiniProDefaultClient()

	ctx := c.Request.Context()

	// 获取用户信息
	userInfo, err := miniprogramClient.Auth.Session(ctx, code)

	if err != nil {
		common.Fail(c, 400, err.Error(), nil)
	}

	fmt.Println(userInfo.OpenID)

	if userInfo.UnionID == "" {
		userInfo.UnionID = userInfo.OpenID
	}

	fmt.Println(userInfo)
}
