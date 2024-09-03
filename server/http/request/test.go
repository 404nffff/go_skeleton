package request

type AuthParams struct {
	Code string `json:"code"  binding:"required"`
}
