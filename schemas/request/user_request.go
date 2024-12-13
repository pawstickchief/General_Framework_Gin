package request

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginUpdatePasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"new_password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
