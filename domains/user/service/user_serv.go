// Package service 定義用戶服務層
package service

import (
	"errors"
	"fmt"
	"self_go_gin/gin_application/api/v1/user/request"

	"self_go_gin/domains/user/entity/model"
	"self_go_gin/domains/user/repository"
	"self_go_gin/gin_application/handler"
	"self_go_gin/util/bcryptencap"
	"self_go_gin/util/jwt_secret"

	"gorm.io/gorm"
)

// UserService 用戶服務層
type UserService struct {
	repo repository.UserRepositoryInterface
}

// NewUserService 創建用戶服務層
func NewUserService() (*UserService, error) {
	repo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("UserService NewUserService(): %w", err)
	}

	return &UserService{
		repo: repo,
	}, nil
}

// CreateUser 創建用戶
func (s *UserService) CreateUser(req request.CreateUserRequest) (*model.User, error) {
	logData := map[string]interface{}{
		"req": req,
	}
	_, err := s.repo.GetUsersByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("UserService CreateUser() data: %s \n %w", logData, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//密碼加密
		bcryptPassword, err := bcryptencap.GenerateFromPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("UserService CreateUser() bcrypt fail data: %s \n %w", logData, err)
		}

		newUsers := &model.User{
			Account:  req.Account,
			Password: string(bcryptPassword),
		}

		user, err := s.repo.CreateUser(newUsers)
		if err != nil {
			return nil, fmt.Errorf("UserService CreateUser() data: %s \n %w", logData, err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("UserService CreateUser() resource exist data: %s \n %w", logData, handler.ErrResourceExist)
}

// CheckLogin 驗證用戶登入
func (s *UserService) CheckLogin(req request.UserLoginRequest) (*string, error) {
	logData := map[string]interface{}{
		"req": req,
	}

	user, err := s.repo.GetUsersByAccount(req.Account)
	if err != nil {
		return nil, fmt.Errorf("UserService CheckLogin() data: %s \n %w", logData, err)
	}

	//密碼驗證
	if err := bcryptencap.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("UserService CheckLogin() CompareHashAndPassword() data : %+v \n %w", logData, err) // 密碼錯誤
	}
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginUser, user.ID)
	if err != nil {
		return nil, fmt.Errorf("UserService CheckLogin() GenerateToken() data : %+v \n %w", logData, err) // 密碼錯誤
	}

	return &jwtToken, nil

}
