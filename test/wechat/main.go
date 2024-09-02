package main

import (
	"context"
	"fmt"
	"tool/pkg/wechat/miniprogram"
)

func main() {

	miniprogramClient := miniprogram.NewMiniProgramClient("Default")
	ctx := context.Background()

	// 获取用户信息
	userInfo, _ := miniprogramClient.Auth.Session(ctx, "0c1OuL100FYnKS11Sk400uJIRl3OuL1i")

	fmt.Println(userInfo)

}
