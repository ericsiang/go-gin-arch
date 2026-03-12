package response

// CreateAdminResponse 創建管理員回應
type CreateAdminResponse struct {
	AdminID uint   `json:"admin_id"`
	Account string `json:"account"`
}
