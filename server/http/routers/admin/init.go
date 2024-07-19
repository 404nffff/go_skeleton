package admin

import (
	"fmt"
	"net/http"
	"tool/global/variable"
	"tool/server/http/middleware"
	"tool/server/http/templates"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var Api *gin.Engine

func InitRouter() *gin.Engine {

	// 使用嵌入的文件系统作为模板文件系统
	tmpl, err := templates.Load()

	if err != nil {
		panic(fmt.Sprintf("Failed to load templates: %v", err))
	}

	//判断是否是调试模式
	AppDebug := variable.ConfigYml.GetBool("AppDebug")

	if AppDebug == true {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		Api = gin.New()
		pprof.Register(Api)
	} else {
		gin.SetMode(gin.ReleaseMode)
		Api = gin.New()
	}

	// 设置模板
	Api.SetHTMLTemplate(tmpl)

	//Api.LoadHTMLGlob(variable.BasePath + "/app/http/templates/**/**/*")

	// 初始化中间件
	initMiddleware()

	//注册静态文件
	//Api.Static("/public/admin", variable.BasePath+"/public/admin")

	staticServer := http.FileServer(http.FS(templates.Components))

	Api.GET("/static/*filepath", func(c *gin.Context) {
		http.StripPrefix("/static/", staticServer).ServeHTTP(c.Writer, c.Request)
	})

	// 注册路由
	RegisterRouter()

	//打印所有嵌入文件
	templates.List()

	return Api
}

// 初始化中间件
func initMiddleware() {

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		Api.Use(middleware.Cors())
	}

	//使用 gin.Recovery() 中间件
	Api.Use(gin.Recovery())

	//使用 LoggerMiddleware 中间件
	logger, err := middleware.InitLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	defer logger.Sync()

	//初始化日志
	Api.Use(middleware.LoggerMiddleware(logger))

	//初始化session
	Api.Use(middleware.SessionMiddleware())

	//初始化参数验证
	Api.Use(middleware.ValidateParams())

	// 使用 panic 恢复中间件
	Api.Use(middleware.PanicRecovery())
}
