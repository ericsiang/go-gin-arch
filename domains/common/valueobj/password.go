// Package valueobj 定義共享的值物件
package valueobj

import (
	"errors"
	"self_go_gin/util/bcryptencap"
)

// Password 密碼值物件
type Password struct {
	hashedValue string
}

// NewPasswordFromPlainText 從明文密碼創建（會自動加密）
// 這是創建新密碼時使用的方法（註冊、修改密碼）
func NewPasswordFromPlainText(plainText string) (Password, error) {
	// 驗證規則 1: 長度檢查
	if plainText == "" {
		return Password{}, errors.New("密碼不能為空")
	}
	if len(plainText) < 6 {
		return Password{}, errors.New("密碼長度至少6位")
	}
	if len(plainText) > 72 {
		// bcrypt 的限制
		return Password{}, errors.New("密碼長度不能超過72位")
	}

	// 可選：密碼強度驗證（可根據需求開啟）
	// if !hasUpperCase(plainText) || !hasLowerCase(plainText) || !hasDigit(plainText) {
	//     return Password{}, errors.New("密碼必須包含大小寫字母和數字")
	// }

	// 加密密碼
	hashed, err := bcryptencap.GenerateFromPassword(plainText)
	if err != nil {
		return Password{}, err
	}

	return Password{hashedValue: string(hashed)}, nil
}

// NewPasswordFromHash 從資料庫讀取的加密密碼創建
// 這是從資料庫載入已存在的用戶時使用的方法
func NewPasswordFromHash(hash string) Password {
	return Password{hashedValue: hash}
}

// Verify 驗證密碼是否正確
// 這是登入時驗證密碼的方法
func (p Password) Verify(plainText string) bool {
	return bcryptencap.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainText)) == nil
}

// Hash 取得加密後的密碼（用於儲存到資料庫）
// 僅供 Repository 層使用
func (p Password) Hash() string {
	return p.hashedValue
}

// Equals 比較兩個密碼值物件是否相等
func (p Password) Equals(other Password) bool {
	return p.hashedValue == other.hashedValue
}

// IsEmpty 檢查密碼是否為空
func (p Password) IsEmpty() bool {
	return p.hashedValue == ""
}
