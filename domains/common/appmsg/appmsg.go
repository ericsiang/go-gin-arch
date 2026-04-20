// Package appmsg 定义所有应用错误消息常量
package appmsg

// 这里集中管理所有的错误消息，方便复用和统一修改
// ===== User Service Messages =====

// User Account Validation
//
//nolint:gosec
const (
	UserAccountFormatInvalid    string = "帳號格式無效"
	UserPasswordInvalidOrWeak   string = "密碼格式或強度不符合要求"
	UserAccountAlreadyExists    string = "帳號已存在"
	UserCheckAccountExistFailed string = "驗證帳號存在狀態失敗"
	UserCreateFailed            string = "建立用戶失敗"
	UserNotFound                string = "用戶不存在"
	UserQueryFailed             string = "查詢用戶失敗"
	UserPasswordIncorrect       string = "帳號或密碼錯誤"
	UserTokenGenerateFailed     string = "生成令牌失敗"
	UserDataCorrupted           string = "用戶數據損壞"
)

// ===== Admin Service Messages =====
//
//nolint:gosec
const (
	AdminAccountFormatInvalid    string = "帳號格式無效"
	AdminPasswordInvalidOrWeak   string = "密碼格式或強度不符合要求"
	AdminAccountAlreadyExists    string = "帳號已存在"
	AdminCheckAccountExistFailed string = "驗證帳號存在狀態失敗"
	AdminCreateFailed            string = "建立管理員失敗"
	AdminNotFound                string = "管理員不存在"
	AdminQueryFailed             string = "查詢管理員失敗"
	AdminPasswordIncorrect       string = "帳號或密碼錯誤"
	AdminTokenGenerateFailed     string = "生成令牌失敗"
	AdminDataCorrupted           string = "管理員數據損壞"
)

// ===== Repository Messages =====
//
//nolint:gosec
const (
	RepositoryInitFailed           string = "初始化倉庫失敗"
	RepositoryQueryFailed          string = "查詢失敗"
	RepositoryCreateFailed         string = "創建失敗"
	RepositoryDataConversionFailed string = "數據轉換失敗"
	DAOInitFailed                  string = "初始化數據訪問層失敗"
	DAODatabaseConnectionFailed    string = "初始化數據庫連接失敗"
	DAOQueryRecordsFailed          string = "查詢記錄失敗"
	DAOCreateRecordFailed          string = "保存記錄失敗"
)

// ===== Event Handler Messages =====
//
//nolint:gosec
const (
	UserEventPayloadParseFailed  string = "事件數據解析失敗"
	AdminEventPayloadParseFailed string = "事件數據解析失敗"
)

// ===== Handler Messages =====
//
//nolint:gosec
const (
	HandlerGetDataFailed    string = "獲取資料失敗"
	HandlerCreateDataFailed string = "創建資料失敗"
	HandlerUpdateDataFailed string = "更新資料失敗"
	HandlerDeleteDataFailed string = "刪除資料失敗"
	HandlerMySQLError       string = "數據庫錯誤"
	HandlerRecordNotFound   string = "記錄不存在"
	HandlerResourceNotFound string = "資源不存在"
	HandlerDeleteNotAllow   string = "刪除不允许"
)
