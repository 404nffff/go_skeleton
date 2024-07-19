package model

type Admin struct {
	ID            int        `gorm:"primaryKey" json:"id"`                      // 主键
	Username      string     `gorm:"type:varchar(20);not null" json:"username"` // 用户名
	Password      string     `gorm:"type:varchar(40);not null" json:"password"` // 密码
	LastLoginTime *LocalTime `gorm:"type:datetime" json:"last_login_time"`      // 上次登录时间
	LoginStatus   int8       `gorm:"type:tinyint" json:"login_status"`          // 登录状态 0禁用 1启用
	CreateTime    *LocalTime `gorm:"type:datetime" json:"create_time"`          // 创建时间
	UpdateTime    *LocalTime `gorm:"type:datetime" json:"update_time"`          // 更新时间
}

// TableName 设置表名前缀
func (Admin) TableName() string {
	return "t_admin"
}
