package api

import (
	"tool/app/http/controller"
	"tool/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func init() {

	//     // 注册路由到 "/api/v1" 路由组
	// registerRoutesToGroup("/api/v1",
	//     route{
	//         method:  "GET",
	//         path:    "/users",
	//         handler: handlers.getUsers,
	//         middlewares: []gin.HandlerFunc{
	//             middlewares.authentication,
	//             middlewares.authorization,
	//         },
	//     },
	//     route{
	//         method:  "POST",
	//         path:    "/users",
	//         handler: handlers.createUser,
	//         middlewares: []gin.HandlerFunc{
	//             middlewares.validation,
	//         },
	//     },
	//     route{
	//         method:  "GET",
	//         path:    "/products",
	//         handler: handlers.getProducts,
	//     },
	// )

	// // 注册中间件到 "/api/v1" 路由组
	// registerMiddlewares("/api/v1",
	//     middlewares.logger,
	//     middlewares.recovery,
	//     middlewares.cors,
	// )

	// // 注册路由到 "/api/v2" 路由组
	// registerRoutesToGroup("/api/v2",
	//     route{
	//         method:  "GET",
	//         path:    "/orders",
	//         handler: handlers.getOrders,
	//     },
	//     route{
	//         method:  "POST",
	//         path:    "/orders",
	//         handler: handlers.createOrder,
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

	registerRoutesToGroup("",
		route{
			method:  "GET",
			path:    "/test",
			handler: controller.Test,
		},
	)

	// 注册路由到 "/user" 路由组
	registerRoutesToGroup("/user",
		route{
			method:  "GET",
			path:    "/test",
			handler: controller.Test,
		},
		route{
			method:  "GET",
			path:    "/test2",
			handler: controller.Test2,
		},
		route{
			method:  "GET",
			path:    "/user",
			handler: controller.GetUsersHandler,
		},
		route{
			method:  "GET",
			path:    "/user2",
			handler: controller.GetUserMongo,
		},
		route{
			method:  "GET",
			path:    "/testudp",
			handler: controller.TestUdp,
		},
		route{
			method:  "GET",
			path:    "/task",
			handler: controller.TestAnt,
		},
		route{
			method:  "GET",
			path:    "/s",
			handler: controller.StatusAnt,
		},
		route{
			method:  "GET",
			path:    "/login",
			handler: controller.GetUsersHandler,
			middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		})

}
