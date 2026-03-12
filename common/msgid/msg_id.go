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
	// NoContent 無內容
	NoContent
	// DuplicateEntry 重複條目
	DuplicateEntry
	// RuleNotAllow 规則不允許
	RuleNotAllow
)
