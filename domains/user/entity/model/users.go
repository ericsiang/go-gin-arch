// Package model 定義用戶相關的數據模型
package model

import (
	"errors"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/internal/model"
)

// User 用戶聚合根
type User struct {
	model.GormModel
	account  valueobj.Account
	password valueobj.Password
}

// NewUser 創建新用戶
func NewUser(account valueobj.Account, password valueobj.Password) *User {
	return &User{
		account:  account,
		password: password,
	}
}

// ReconstructUser 從資料庫重建用戶（用於 Repository 層）
func ReconstructUser(id uint, account valueobj.Account, password valueobj.Password, gormModel model.GormModel) *User {
	user := &User{
		GormModel: gormModel,
		account:   account,
		password:  password,
	}
	user.ID = id
	return user
}

// ============ 業務方法（領域邏輯） ============

// ChangePassword 修改密碼
// 業務規則：
// 1. 必須驗證舊密碼正確
// 2. 新密碼不能與舊密碼相同
func (u *User) ChangePassword(oldPasswordPlain, newPasswordPlain string) error {
	// 驗證舊密碼
	if !u.password.Verify(oldPasswordPlain) {
		return errors.New("舊密碼錯誤")
	}

	// 檢查新舊密碼是否相同
	if oldPasswordPlain == newPasswordPlain {
		return errors.New("新密碼不能與舊密碼相同")
	}

	// 創建新密碼值物件（自動驗證和加密）
	newPassword, err := valueobj.NewPasswordFromPlainText(newPasswordPlain)
	if err != nil {
		return err
	}

	u.password = newPassword
	return nil
}

// VerifyPassword 驗證密碼是否正確
// 用於登入驗證
func (u *User) VerifyPassword(plainText string) bool {
	return u.password.Verify(plainText)
}

// ChangeAccount 修改帳號
// 業務規則：新帳號必須符合格式要求
func (u *User) ChangeAccount(newAccount valueobj.Account) error {
	if u.account.Equals(newAccount) {
		return errors.New("新帳號與舊帳號相同")
	}
	u.account = newAccount
	return nil
}

// ============ 查詢方法（Getter） ============

// GetAccount 取得帳號值
func (u *User) GetAccount() string {
	return u.account.Value()
}

// GetAccountvalueobj 取得帳號值物件
func (u *User) GetAccountvalueobj() valueobj.Account {
	return u.account
}

// GetPasswordHash 取得加密後的密碼（僅供 Repository 層使用）
func (u *User) GetPasswordHash() string {
	return u.password.Hash()
}
