// Package service 定義管理員服務層
package service

import (
	"errors"
	"fmt"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/admin/entity"
	"self_go_gin/domains/admin/repository"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/gin_application/api/v1/admin/request"
	apperror "self_go_gin/internal/apperror"
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
func (s *AdminService) CreateAdmin(req request.CreateAdminRequest) (*entity.Admin, error) {
	// 創建帳號值物件（自動驗證格式）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AdminAccountFormatInvalid,
			fmt.Errorf("invalid account format: %w", err),
			map[string]interface{}{
				"account": req.Account,
			},
		)
	}

	// 創建密碼值物件（自動驗證強度和加密）
	password, err := valueobj.NewPasswordFromPlainText(req.Password)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AdminPasswordInvalidOrWeak,
			fmt.Errorf("invalid password: %w", err),
			nil,
		)
	}

	// 檢查帳號是否已存在
	_, err = s.repo.GetAdminByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.AdminCheckAccountExistFailed,
			fmt.Errorf("check account existence failed: %w", err),
			map[string]interface{}{
				"account": req.Account,
			},
		)
	}
	if err == nil {
		// 帳號已存在
		return nil, apperror.NewAppError(
			msgid.ResourceExist,
			appmsg.AdminAccountAlreadyExists,
			fmt.Errorf("account already exists"),
			nil,
		)
	}

	// 創建聚合根
	admin := entity.NewAdmin(account, password)

	// 儲存到資料庫
	createdAdmin, err := s.repo.CreateAdmin(admin)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.AdminCreateFailed,
			fmt.Errorf("create admin failed: %w", err),
			map[string]interface{}{
				"admin": admin,
			},
		)
	}

	return createdAdmin, nil
}

// CheckLogin 驗證管理員登入
func (s *AdminService) CheckLogin(req request.AdminLoginRequest) (*string, error) {
	// 先驗證帳號格式（快速失敗）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AdminAccountFormatInvalid,
			fmt.Errorf("invalid account format: %w", err),
			map[string]interface{}{
				"account": req.Account,
			},
		)
	}

	// 查詢管理員
	admin, err := s.repo.GetAdminByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewAppError(
				msgid.NoContent,
				appmsg.AdminNotFound,
				fmt.Errorf("admin not found"),
				map[string]interface{}{
					"admin": admin,
				},
			)
		}
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.AdminQueryFailed,
			fmt.Errorf("get admin failed: %w", err),
			map[string]interface{}{
				"admin": admin,
			},
		)
	}

	// 驗證密碼（業務邏輯在聚合根中）
	if !admin.VerifyPassword(req.Password) {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AdminPasswordIncorrect,
			fmt.Errorf("password incorrect"),
			nil,
		)
	}

	// 生成 JWT Token
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginAdmin, admin.ID)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.AdminTokenGenerateFailed,
			fmt.Errorf("generate token failed: %w", err),
			nil,
		)
	}

	return &jwtToken, nil
}
