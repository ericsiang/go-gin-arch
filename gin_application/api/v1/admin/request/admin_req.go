package request

// CreateAdminRequest 創建管理員請求
type CreateAdminRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginRequest 管理員登入請求
type AdminLoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GetAdminsByIDRequest 根據ID獲取管理員請求
type GetAdminsByIDRequest struct {
	FilterAdminsID string `form:"filterAdminsId" json:"filterAdminsId" binding:"required"`
}
