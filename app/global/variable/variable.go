package variable

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"tool/pkg/ants"
	"tool/pkg/yml_config/ymlconfig_interf"

	"go.uber.org/zap"
)

// 初始化参数
var (
	BasePath string // 定义项目的根目录

	EventDestroyPrefix = "Destroy_" //  程序退出时需要销毁的事件前缀

	ConfigYml ymlconfig_interf.YmlConfigInterf // 全局配置文件指针

	Logs *zap.Logger // 全局日志指针

	Pool ants.AntsInterface // 全局协程池指针
)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}

		//判断是否在 D:\\www\\BaiduSyncdisk\\www\\go\\tool\\cmd\\api\\config 目录下 如果在 则切换到 D:\\www\\BaiduSyncdisk\\www\\go\\tool\\

		// 根据操作系统类型替换路径部分
		if runtime.GOOS == "windows" {
			// Windows 路径替换
			BasePath = strings.Replace(BasePath, "cmd\\api", "", -1)
		} else {
			// Linux 路径替换
			BasePath = strings.Replace(BasePath, "cmd/api", "", -1)
		}

		// 清理路径，移除多余的斜杠
		BasePath = strings.TrimSuffix(BasePath, "/")
		BasePath = strings.TrimSuffix(BasePath, "\\")

		fmt.Println("BasePath:", BasePath)

	} else {
		log.Fatal("BasePath error")
	}
}
