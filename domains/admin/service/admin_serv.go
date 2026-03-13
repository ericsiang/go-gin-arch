// Package service 定義管理員服務層
package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"self_go_gin/domains/admin/entity/model"
	"self_go_gin/domains/admin/repository"
	"self_go_gin/gin_application/api/v1/admin/request"
	"self_go_gin/gin_application/handler"
	"self_go_gin/util/bcryptencap"
	"self_go_gin/util/jwt_secret"
)

// AdminService 管理員服務層
type AdminService struct {
	repo repository.AdminRepositoryInterface
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
	logData := map[string]interface{}{
		"req": req,
	}
	_, err := s.repo.GetAdminByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("AdminService CreateAdmin() data: %s \n %w", logData, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//密碼加密
		bcryptPassword, err := bcryptencap.GenerateFromPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("AdminService CreateAdmin() bcrypt fail data: %s \n %w", logData, err)
		}

		newAdmins := model.Admins{
			Account:  req.Account,
			Password: string(bcryptPassword),
		}

		user, err := s.repo.CreateAdmin(newAdmins)
		if err != nil {
			return nil, fmt.Errorf("AdminService AdminService() data: %s \n %w", logData, err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("AdminService AdminService() resource exist data: %s \n %w", logData, handler.ErrResourceExist)
}

// CheckLogin 驗證管理員登入
func (s *AdminService) CheckLogin(req request.AdminLoginRequest) (*string, error) {
	logData := map[string]interface{}{
		"req": req,
	}
	admin, err := s.repo.GetAdminByAccount(req.Account)
	if err != nil {
		return nil, fmt.Errorf("AdminService CheckLogin() data: %s \n %w", logData, err)
	}

	//密碼驗證
	if err := bcryptencap.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("AdminService CheckLogin() CompareHashAndPassword() data : %+v \n %w", logData, err) // 密碼錯誤
	}

	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginAdmin, admin.ID)
	if err != nil {
		return nil, fmt.Errorf("AdminService CheckLogin() GenerateToken() data : %+v \n %w", logData, err) // 密碼錯誤
	}

	return &jwtToken, nil

}
