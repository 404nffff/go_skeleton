package api

import "reflect"

// 定义不同路由所需的请求参数结构体
type ParamsForRoute1 struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required,email"`
}

type ParamsForRoute2 struct {
	Age int `form:"age" binding:"required,gte=1,lte=120"`
}

// 路径到结构体类型的映射
var RouteParamMap = map[string]reflect.Type{
	"/user/user2": reflect.TypeOf(ParamsForRoute1{}),
	"/route2": reflect.TypeOf(ParamsForRoute2{}),
}
