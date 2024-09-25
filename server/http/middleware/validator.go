package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"tool/global/utils/common"
	"tool/pkg/web_server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 定义全局验证器
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name == "" {
			name = fld.Name
		}
		return name
	})
}

// 通用的参数验证中间件生成器
func ValidateParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		paramType, exists := web_server.RouteParamMap[route]
		if !exists {
			c.Set("params", nil)
			c.Next()
			return
		}

		params := reflect.New(paramType).Interface()

		// 根据 Content-Type 绑定参数
		contentType := c.GetHeader("Content-Type")
		
		switch {
		case strings.HasPrefix(contentType, "application/json"):
			if err := c.ShouldBindJSON(params); err != nil {
				handleValidationError(c, err)
				return
			}
		case strings.HasPrefix(contentType, "application/xml"):
			if err := c.ShouldBindXML(params); err != nil {
				handleValidationError(c, err)
				return
			}
		case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"), strings.HasPrefix(contentType, "multipart/form-data"):
			if err := c.ShouldBind(params); err != nil {
				handleValidationError(c, err)
				return
			}
		case c.Request.Method == http.MethodGet:
			if err := c.ShouldBindQuery(params); err != nil {
				handleValidationError(c, err)
				return
			}
		default:
			handleValidationError(c, errors.New("unsupported content type"))
			return
		}


		// 验证参数
		if err := validate.Struct(params); err != nil {
			handleValidationError(c, err)
			return
		}

		// 传递参数
		c.Set("params", params)

		c.Next()
	}
}

// 错误处理函数
func handleValidationError(c *gin.Context, err error) {
	var verr validator.ValidationErrors
	var errorMessages []string
	if errors.As(err, &verr) {
		for _, fieldErr := range verr {
			errorMessages = append(errorMessages, getErrorMessage(fieldErr))
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	common.Fail(c, http.StatusBadRequest, strings.Join(errorMessages, "; "), nil)
}

// 获取自定义的错误消息
func getErrorMessage(fe validator.FieldError) string {
	fieldName := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s 是必填项", fieldName)
	case "email":
		return fmt.Sprintf("%s 必须是一个有效的电子邮件地址", fieldName)
	case "gte":
		return fmt.Sprintf("%s 必须大于或等于 %s", fieldName, fe.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于或等于 %s", fieldName, fe.Param())
	case "min":
		return fmt.Sprintf("%s 的最小值是 %s", fieldName, fe.Param())
	case "max":
		return fmt.Sprintf("%s 的最大值是 %s", fieldName, fe.Param())
	case "len":
		return fmt.Sprintf("%s 的长度必须是 %s", fieldName, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s 必须是 [%s] 其中之一", fieldName, fe.Param())
	case "url":
		return fmt.Sprintf("%s 必须是一个有效的URL", fieldName)
	case "uuid":
		return fmt.Sprintf("%s 必须是一个有效的UUID", fieldName)
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", fieldName)
	case "numeric":
		return fmt.Sprintf("%s 必须是一个数字", fieldName)
	case "boolean":
		return fmt.Sprintf("%s 必须是布尔值", fieldName)
	case "datetime":
		return fmt.Sprintf("%s 必须是一个有效的日期时间格式", fieldName)
	default:
		return fmt.Sprintf("%s 字段验证错误", fieldName)
	}
}
