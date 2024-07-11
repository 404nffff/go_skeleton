package db_client

import (
	"tool/pkg/mysql"

	"gorm.io/gorm"
)

func MysqlLocal() *gorm.DB {

	return mysql.NewClient("Local")
}
