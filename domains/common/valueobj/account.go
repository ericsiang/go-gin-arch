// Package valueobj 定義共享的值物件
package valueobj

import (
	"errors"
	"regexp"
	"strings"
)

// Account 帳號值物件
type Account struct {
	value string
}

// NewAccount 創建帳號值物件，包含驗證邏輯
func NewAccount(account string) (Account, error) {
	// 移除前後空白
	account = strings.TrimSpace(account)

	// 驗證規則 1: 不能為空
	if account == "" {
		return Account{}, errors.New("帳號不能為空")
	}

	// 驗證規則 2: 長度限制
	if len(account) < 3 {
		return Account{}, errors.New("帳號長度至少3個字元")
	}
	if len(account) > 255 {
		return Account{}, errors.New("帳號長度不能超過255個字元")
	}

	// 驗證規則 3: 格式驗證（只允許字母、數字、底線、減號、點）
	// 可根據實際需求調整
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !validPattern.MatchString(account) {
		return Account{}, errors.New("帳號只能包含字母、數字、底線、減號和點")
	}

	// 驗證規則 4: 不能以特殊字元開頭或結尾
	if strings.HasPrefix(account, ".") || strings.HasSuffix(account, ".") ||
		strings.HasPrefix(account, "-") || strings.HasSuffix(account, "-") {
		return Account{}, errors.New("帳號不能以點或減號開頭或結尾")
	}

	return Account{value: account}, nil
}

// Value 取得帳號值
func (a Account) Value() string {
	return a.value
}

// Equals 比較兩個帳號值物件是否相等（不區分大小寫）
func (a Account) Equals(other Account) bool {
	return strings.EqualFold(a.value, other.value)
}

// EqualsString 比較帳號值物件與字串是否相等（不區分大小寫）
func (a Account) EqualsString(value string) bool {
	return strings.EqualFold(a.value, value)
}

// IsEmpty 檢查帳號是否為空
func (a Account) IsEmpty() bool {
	return a.value == ""
}

// String 實作 Stringer 介面，方便日誌輸出
func (a Account) String() string {
	return a.value
}
