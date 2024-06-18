package validator

import (
	"reflect"
	"tool/app/http/validator/admin"
)

// 路径到结构体类型的映射
var RouteParamMap = map[string]reflect.Type{
	"/admin/login_submit": reflect.TypeOf(admin.LoginValidate{}),
}
