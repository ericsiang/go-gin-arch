// Package response 定義管理員相關的API回應結構
package response

// CreateAdminResponse 創建管理員回應
type CreateAdminResponse struct {
	AdminID uint   `json:"admin_id"`
	Account string `json:"account"`
}
