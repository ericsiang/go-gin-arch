// Package appmsg 定义所有应用错误消息常量
package appmsg

// 这里集中管理所有的错误消息，方便复用和统一修改

// Account Validation
//
//nolint:gosec
const (
	AccountFormatInvalid    string = "帳號格式無效"
	PasswordInvalidOrWeak   string = "密碼格式或強度不符合要求"
	AccountAlreadyExists    string = "帳號已存在"
	CheckAccountExistFailed string = "驗證帳號存在狀態失敗"
	PasswordIncorrect       string = "帳號或密碼錯誤"
	TokenGenerateFailed     string = "生成令牌失敗"
)

// ===== DB Messages =====
//
//nolint:gosec
const (
	InitFailed               string = "初始化失敗"
	QueryFailed              string = "查詢失敗"
	CreateFailed             string = "創建失敗"
	DataConversionFailed     string = "數據轉換失敗"
	DatabaseConnectionFailed string = "初始化數據庫連接失敗"
	RecordNotFound           string = "記錄不存在"
	DuplicateEntry           string = "資源重複"
)

// ===== Event Handler Messages =====
//
//nolint:gosec
const (
	UserEventPayloadParseFailed  string = "事件數據解析失敗"
	UserEventPublishFailed       string = "事件發布失敗"
	AdminEventPayloadParseFailed string = "事件數據解析失敗"
)


// ===== Redis Messages =====
//
//nolint:gosec
const (
	RedisSetFailed string = "Redis Set 操作失敗"
)