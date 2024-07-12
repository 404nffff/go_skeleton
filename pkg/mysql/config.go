package mysql

import (
	"fmt"
	"tool/app/global/variable"
	"tool/pkg/yml_config"
)

// DatabaseConfig 定义数据库配置结构体
type DatabaseConfig struct {
	User               string // 数据库用户名
	Pass               string // 数据库密码
	Host               string // 数据库地址
	Port               string // 数据库端口
	Database           string // 数据库名称
	Charset            string // 数据库字符集
	SetMaxIdleConns    int    // 连接池中的最大空闲连接数
	SetMaxOpenConns    int    // 数据库的最大连接数量
	SetConnMaxLifetime int    // 连接的最大可复用时间
	EventDestroyPrefix string // 事件销毁前缀
}

// 加载配置文件
func loadConfig(conn string) DatabaseConfig {

	mysqlConfig := yml_config.LoadConfig("mysql")

	//查找配置文件中的数据库配置
	// Local:
	// 	Host: "127.0.0.1"
	// 	DataBase: ""
	// 	Port:
	// 	User: ""
	// 	Pass: ""
	// 	Charset: "utf8"
	// 	SetMaxIdleConns: 10
	// 	SetMaxOpenConns: 128
	// 	SetConnMaxLifetime: 60    # 连接不活动时的最大生存时间(秒)
	// 	SlowThreshold: 30            # 慢 SQL 阈值(sql执行时间超过此时间单位（秒），就会触发系统日志记录)

	if mysqlConfig.GetString(conn+".User") == "" {
		panic(fmt.Sprintf("Failed to get Mysql config: %s", conn))
	}

	config := DatabaseConfig{
		User:               mysqlConfig.GetString(conn + ".User"),
		Pass:               mysqlConfig.GetString(conn + ".Pass"),
		Host:               mysqlConfig.GetString(conn + ".Host"),
		Port:               mysqlConfig.GetString(conn + ".Port"),
		Database:           mysqlConfig.GetString(conn + ".DataBase"),
		Charset:            mysqlConfig.GetString(conn + ".Charset"),
		SetMaxIdleConns:    mysqlConfig.GetInt(conn + ".SetMaxIdleConns"),
		SetMaxOpenConns:    mysqlConfig.GetInt(conn + ".SetMaxOpenConns"),
		SetConnMaxLifetime: mysqlConfig.GetInt(conn + ".SetConnMaxLifetime"),
		EventDestroyPrefix: variable.EventDestroyPrefix + "Mysql_" + conn,
	}

	return config
}
