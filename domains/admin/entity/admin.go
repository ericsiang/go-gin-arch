// Package entity 定義管理員相關的數據模型
package entity

import (
	"errors"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/internal/model"
)

// Admin 管理員聚合根
type Admin struct {
	model.GormModel
	account  valueobj.Account
	password valueobj.Password
}

// NewAdmin 創建新管理員
func NewAdmin(account valueobj.Account, password valueobj.Password) *Admin {
	return &Admin{
		account:  account,
		password: password,
	}
}

// ReconstructAdmin 從資料庫重建管理員（用於 Repository 層）
// 包含完整的 GORM 模型資料
func ReconstructAdmin(id uint, account valueobj.Account, password valueobj.Password, gormModel model.GormModel) *Admin {
	admin := &Admin{
		GormModel: gormModel,
		account:   account,
		password:  password,
	}
	admin.ID = id
	return admin
}

// ============ 業務方法（領域邏輯） ============

// ChangePassword 修改密碼
// 業務規則：
// 1. 必須驗證舊密碼正確
// 2. 新密碼不能與舊密碼相同
func (a *Admin) ChangePassword(oldPasswordPlain, newPasswordPlain string) error {
	// 驗證舊密碼
	if !a.password.Verify(oldPasswordPlain) {
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

	a.password = newPassword
	return nil
}

// VerifyPassword 驗證密碼是否正確
// 用於登入驗證
func (a *Admin) VerifyPassword(plainText string) bool {
	return a.password.Verify(plainText)
}

// ChangeAccount 修改帳號
// 業務規則：新帳號必須符合格式要求
func (a *Admin) ChangeAccount(newAccount valueobj.Account) error {
	if a.account.Equals(newAccount) {
		return errors.New("新帳號與舊帳號相同")
	}
	a.account = newAccount
	return nil
}

// ============ 查詢方法（Getter） ============

// GetAccount 取得帳號值
func (a *Admin) GetAccount() string {
	return a.account.Value()
}

// GetAccountValueObject 取得帳號值物件
func (a *Admin) GetAccountValueObject() valueobj.Account {
	return a.account
}

// GetPasswordHash 取得加密後的密碼（僅供 Repository 層使用）
func (a *Admin) GetPasswordHash() string {
	return a.password.Hash()
}
