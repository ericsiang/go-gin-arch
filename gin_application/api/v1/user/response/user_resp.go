package response

// CreateUserResponse 創建用戶回應
type CreateUserResponse struct {
	UsersID uint   `json:"id"`
	Account string `json:"account"`
}
