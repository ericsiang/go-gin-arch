// Package response 定義用戶相關的API回應結構
package response

// CreateUserResponse 創建用戶回應
type CreateUserResponse struct {
	UsersID uint   `json:"id"`
	Account string `json:"account"`
}
