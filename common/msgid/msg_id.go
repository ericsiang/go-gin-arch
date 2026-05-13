// Package msgid 定義訊息識別碼
package msgid

// MsgID response 訊息的識別碼
type MsgID uint32

const (
	// Success 成功
	Success MsgID = iota
	// Fail 失敗
	Fail
	// TokenExpires 令牌過期
	TokenExpires
	// TokenInvalid 令牌無效
	TokenInvalid
	// LoginNotExist 登錄不存在
	LoginNotExist
	// NoContent 無內容
	NoContent
	// DuplicateEntry 重複條目
	DuplicateEntry
	// RuleNotAllow 規則不允許
	RuleNotAllow
	// ResourceExist 資源已存在
	ResourceExist
	// InvalidInput 輸入驗證失敗（帳號格式、密碼強度等）
	InvalidInput
	// NotAppError 非 AppError 包裝過的錯誤
	NotAppError
)
