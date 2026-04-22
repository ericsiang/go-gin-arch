// Package service 定義管理員服務層
package service

import (
	"errors"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/admin/entity"
	"self_go_gin/domains/admin/repository"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/gin_application/api/v1/admin/request"
	"self_go_gin/gin_application/handler"
	apperror "self_go_gin/internal/apperror"
	jwtsecret "self_go_gin/util/jwt_secret"

	"gorm.io/gorm"
)

// AdminService 管理員服務層
type AdminService struct {
	repo repository.AdminRepository
}

// NewAdminService 創建管理員服務層
func NewAdminService() (*AdminService, *apperror.AppError) {
	repo, err := repository.NewAdminRepository()
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.InitFailed,
			err,
			apperror.WithLayer("AdminService NewAdminService()"),
		)
	}
	return &AdminService{
		repo: repo,
	}, nil
}

// CreateAdmin 創建管理員
func (s *AdminService) CreateAdmin(req request.CreateAdminRequest) (*entity.Admin, *apperror.AppError) {
	// 創建帳號值物件（自動驗證格式）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AccountFormatInvalid,
			err,
			apperror.WithLayer("AdminService CreateAdmin()"),
			apperror.WithLogData(map[string]interface{}{
				"account": req.Account,
			}),
		)
	}

	// 創建密碼值物件（自動驗證強度和加密）
	password, err := valueobj.NewPasswordFromPlainText(req.Password)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.PasswordInvalidOrWeak,
			err,
			apperror.WithLayer("AdminService CreateAdmin() NewPasswordFromPlainText()"),
		)
	}

	// 檢查帳號是否已存在
	_, err = s.repo.GetAdminByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CheckAccountExistFailed,
			err,
			apperror.WithLayer("AdminService CreateAdmin() GetAdminByAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": req.Account,
			}),
		)
	}
	if err == nil {
		// 帳號已存在
		return nil, apperror.NewAppError(
			msgid.ResourceExist,
			appmsg.AccountAlreadyExists,
			handler.ErrResourceExist,
			apperror.WithLayer("AdminService CreateAdmin() GetAdminByAccount() ResourceExist"),
		)
	}

	// 創建聚合根
	admin := entity.NewAdmin(account, password)

	// 儲存到資料庫
	createdAdmin, err := s.repo.CreateAdmin(admin)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CreateFailed,
			err,
			apperror.WithLayer("AdminService CreateAdmin() CreateAdmin()"),
			apperror.WithLogData(map[string]interface{}{
				"admin": admin,
			}),
		)
	}

	return createdAdmin, nil
}

// CheckLogin 驗證管理員登入
func (s *AdminService) CheckLogin(req request.AdminLoginRequest) (*string, *apperror.AppError) {
	// 先驗證帳號格式（快速失敗）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AccountFormatInvalid,
			err,
			apperror.WithLayer("AdminService CheckLogin() NewAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": req.Account,
			}),
		)
	}

	// 查詢管理員
	admin, err := s.repo.GetAdminByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewAppError(
				msgid.NoContent,
				appmsg.RecordNotFound,
				handler.ErrRecordNotFound,
				apperror.WithLayer("AdminService CheckLogin() GetAdminByAccount() RecordNotFound"),
				apperror.WithLogData(map[string]interface{}{
					"admin": admin,
				}),
			)
		}
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.QueryFailed,
			err,
			apperror.WithLayer("AdminService CheckLogin() GetAdminByAccount() QueryFailed"),
			apperror.WithLogData(map[string]interface{}{
				"admin": admin,
			}),
		)
	}

	// 驗證密碼（業務邏輯在聚合根中）
	if !admin.VerifyPassword(req.Password) {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.PasswordIncorrect,
			handler.ErrPasswordIncorrect,
			apperror.WithLayer("AdminService CheckLogin() VerifyPassword()"),
		)
	}

	// 生成 JWT Token
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginAdmin, admin.ID)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.TokenGenerateFailed,
			err,
			apperror.WithLayer("AdminService CheckLogin() GenerateToken()"),
		)
	}

	return &jwtToken, nil
}
