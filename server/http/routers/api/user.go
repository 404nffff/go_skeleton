package api

import (
	"tool/pkg/web_server"
	"tool/server/http/controller"
	"tool/server/http/middleware"

	"github.com/gin-gonic/gin"
)

func init() {

	//     // 注册路由到 "/api/v1" 路由组
	// web_server.RegisterRoutes("/api/v1",
	//     route{
	//         Method:  "GET",
	//         Path:    "/users",
	//         Handler: Handlers.getUsers,
	//         middlewares: []gin.HandlerFunc{
	//             middlewares.authentication,
	//             middlewares.authorization,
	//         },
	//     },
	//     route{
	//         Method:  "POST",
	//         Path:    "/users",
	//         Handler: Handlers.createUser,
	//         middlewares: []gin.HandlerFunc{
	//             middlewares.validation,
	//         },
	//     },
	//     route{
	//         Method:  "GET",
	//         Path:    "/products",
	//         Handler: Handlers.getProducts,
	//     },
	// )

	// // 注册中间件到 "/api/v1" 路由组
	// registerMiddlewares("/api/v1",
	//     middlewares.logger,
	//     middlewares.recovery,
	//     middlewares.cors,
	// )

	// // 注册路由到 "/api/v2" 路由组
	// web_server.RegisterRoutes("/api/v2",
	//     route{
	//         Method:  "GET",
	//         Path:    "/orders",
	//         Handler: Handlers.getOrders,
	//     },
	//     route{
	//         Method:  "POST",
	//         Path:    "/orders",
	//         Handler: Handlers.createOrder,
	//         middlewares: []gin.HandlerFunc{
	//             middlewares.authentication,
	//         },
	//     },
	// )

	// // 注册中间件到 "/api/v2" 路由组
	// registerMiddlewares("/api/v2",
	//     middlewares.logger,
	//     middlewares.recovery,
	// )

	//web_server.RegisterMiddleware("/user", middleware.AuthMiddleware())

	web_server.RegisterRoutes("",
		web_server.Route{
			Method:  "GET",
			Path:    "/test",
			Handler: controller.Test,
		},
	)

	// 注册路由到 "/user" 路由组
	web_server.RegisterRoutes("/user",
		web_server.Route{
			Method:  "GET",
			Path:    "/test",
			Handler: controller.Test,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/test2",
			Handler: controller.Test2,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/user",
			Handler: controller.GetUsersHandler,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/user2",
			Handler: controller.GetUserMongo,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/testudp",
			Handler: controller.TestUdp,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/task",
			Handler: controller.TestAnt,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/s",
			Handler: controller.StatusAnt,
		},
		web_server.Route{
			Method:  "GET",
			Path:    "/login",
			Handler: controller.GetUsersHandler,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		})

}
