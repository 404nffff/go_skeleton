package admin

type LoginValidate struct {
	Username string `form:"username" binding:"required,alphanum"`
	Password string `form:"password" binding:"required,alphanum"`
}
