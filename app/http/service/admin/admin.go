package admin

import (
	"fmt"
	"tool/app/http/model"
	"tool/app/utils/common"
	"tool/app/utils/db_client"
)

var mysql = db_client.MysqlLocal()

// Login 登录函数
func Login(data map[string]any) map[string]any {
	username := data["username"].(string)
	// password := data["password"].(string)

	var users model.Admin
	if err := mysql.Where(&model.Admin{Username: username}).Find(&users).Error; err != nil {
		return common.ServiceResponse(400, "用户或密码错误", nil)
	}

	fmt.Println(common.Md5(users.Password))
	fmt.Println(data["password"])

	//md5 判断
	if users.Password != common.Md5(data["password"].(string)) {
		return common.ServiceResponse(400, "用户或密码错误", nil)
	}

	returnData := map[string]any{
		"id":       users.ID,
		"username": users.Username,
	}

	return common.ServiceResponse(200, "登录成功", returnData)
}
