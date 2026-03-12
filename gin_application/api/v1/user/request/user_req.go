package request

// CreateUserRequest 創建用戶請求
type CreateUserRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginRequest 用戶登入請求
type UserLoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GetUsersByIDRequest 根據ID獲取用戶請求
type GetUsersByIDRequest struct {
	FilterUsersID string `form:"filterUsersID" json:"filterUsersID" binding:"required"`
}
