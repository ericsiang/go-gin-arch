// Package service 定義管理員服務層
package service

import (
	"errors"
	"fmt"
	"self_go_gin/domains/admin/entity/model"
	"self_go_gin/domains/admin/repository"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/gin_application/api/v1/admin/request"
	"self_go_gin/gin_application/handler"
	jwtsecret "self_go_gin/util/jwt_secret"

	"gorm.io/gorm"
)

// AdminService 管理員服務層
type AdminService struct {
	repo repository.AdminRepository
}

// NewAdminService 創建管理員服務層
func NewAdminService() (*AdminService, error) {
	repo, err := repository.NewAdminRepository()
	if err != nil {
		return nil, fmt.Errorf("AdminService NewAdminService() : %w", err)
	}
	return &AdminService{
		repo: repo,
	}, nil
}

// CreateAdmin 創建管理員
func (s *AdminService) CreateAdmin(req request.CreateAdminRequest) (*model.Admins, error) {
	// 創建帳號值物件（自動驗證格式）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, fmt.Errorf("invalid account: %w", err)
	}

	// 創建密碼值物件（自動驗證強度和加密）
	password, err := valueobj.NewPasswordFromPlainText(req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// 檢查帳號是否已存在
	_, err = s.repo.GetAdminByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check account existence failed: %w", err)
	}
	if err == nil {
		// 帳號已存在
		return nil, fmt.Errorf("account already exists: %w", handler.ErrResourceExist)
	}

	// 創建聚合根
	admin := model.NewAdmins(account, password)

	// 儲存到資料庫
	createdAdmin, err := s.repo.CreateAdmin(admin)
	if err != nil {
		return nil, fmt.Errorf("create admin failed: %w", err)
	}

	return createdAdmin, nil
}

// CheckLogin 驗證管理員登入
func (s *AdminService) CheckLogin(req request.AdminLoginRequest) (*string, error) {
	// 先驗證帳號格式（快速失敗）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, fmt.Errorf("invalid account format: %w", err)
	}

	// 查詢管理員
	admin, err := s.repo.GetAdminByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("get admin failed: %w", err)
	}

	// 驗證密碼（業務邏輯在聚合根中）
	if !admin.VerifyPassword(req.Password) {
		return nil, fmt.Errorf("password incorrect")
	}

	// 生成 JWT Token
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginAdmin, admin.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	return &jwtToken, nil
}
