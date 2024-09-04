package wechat

import (
	"tool/pkg/wechat/miniprogram"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

// MiniProDefaultClient 获取默认小程序客户端
func MiniProDefaultClient() *miniProgram.MiniProgram {
	miniprogramClient := miniprogram.NewMiniProgramClient("Default")
	return miniprogramClient
}
