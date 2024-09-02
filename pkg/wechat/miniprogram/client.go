package miniprogram

import (
	"fmt"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

func NewMiniProgramClient(name string) *miniProgram.MiniProgram {

	config := loadConfig(name)

	MiniProgramApp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     config.AppID,  // 小程序appid
		Secret:    config.Secret, // 小程序app secret
		HttpDebug: config.Debug,
		Log: miniProgram.Log{
			Level: config.Log.Level,
			File:  config.Log.Path,
		},
		// 可选，不传默认走程序内存
		Cache: kernel.NewRedisClient(&kernel.UniversalOptions{
			Addrs:    []string{config.Redis.Addr},
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		}),
	})

	fmt.Println(MiniProgramApp)

	if err != nil {
		return nil
	}

	return MiniProgramApp
}
