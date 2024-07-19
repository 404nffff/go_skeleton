package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"tool/global/utils/common"
	web_validate "tool/server/http/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 定义全局验证器
var validate = validator.New()

// 通用的参数验证中间件生成器
func ValidateParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		paramType, exists := web_validate.RouteParamMap[route]
		if !exists {
			c.Next()
			return
		}

		params := reflect.New(paramType).Interface()

		// 绑定参数
		if err := c.ShouldBind(params); err != nil {
			handleValidationError(c, err)
			return
		}

		// 验证参数
		if err := validate.Struct(params); err != nil {
			handleValidationError(c, err)
			return
		}
		c.Next()
	}
}

// 错误处理函数
func handleValidationError(c *gin.Context, err error) {
	var verr validator.ValidationErrors
	var errorMessages string
	if errors.As(err, &verr) {

		for _, fieldErr := range verr {
			errorMessages = getErrorMessage(fieldErr)
			break
		}
	} else {
		errorMessages = err.Error()
	}

	common.Fail(c, http.StatusBadRequest, errorMessages, nil)
}

// 获取自定义的错误消息
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s 是必填项", fe.Field())
	case "email":
		return fmt.Sprintf("%s 必须是一个有效的电子邮件地址", fe.Field())
	case "gte":
		return fmt.Sprintf("%s 必须大于或等于 %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于或等于 %s", fe.Field(), fe.Param())
	}
	return fmt.Sprintf("%s 字段验证错误", fe.Field())
}
