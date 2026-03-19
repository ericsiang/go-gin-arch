// Package service 定義用戶服務層
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"self_go_gin/container"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/domains/user/entity/model"
	"self_go_gin/domains/user/events"
	"self_go_gin/domains/user/repository"
	"self_go_gin/gin_application/api/v1/user/request"
	"self_go_gin/gin_application/handler"
	"self_go_gin/infra/event"
	jwtsecret "self_go_gin/util/jwt_secret"

	"gorm.io/gorm"
)

// UserService 用戶服務層
type UserService struct {
	repo      repository.UserRepository
	publisher event.Publisher
}

// NewUserService 創建用戶服務層
func NewUserService() (*UserService, error) {
	repo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("UserService NewUserService(): %w", err)
	}
	app := container.GetContainer()
	if app.GetConfig().IsEventBroker {
		broker := app.GetEventBroker()
		return &UserService{
			repo:      repo,
			publisher: broker.Publisher(), // 使用工廠獲取事件發布器
		}, nil
	}else{
		return &UserService{
			repo:      repo,
		}, nil
	}

}

// CreateUser 創建用戶
func (s *UserService) CreateUser(req request.CreateUserRequest) (*model.User, error) {
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
	_, err = s.repo.GetUsersByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check account existence failed: %w", err)
	}
	if err == nil {
		// 帳號已存在
		return nil, fmt.Errorf("account already exists: %w", handler.ErrResourceExist)
	}

	// 創建聚合根
	user := model.NewUser(account, password)

	// 儲存到資料庫
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	// 發布用戶創建事件
	if s.publisher != nil {
		if err := s.publishUserCreatedEvent(context.Background(), createdUser); err != nil {
			// 記錄錯誤但不阻止用戶創建流程
			log.Printf("[UserService] Failed to publish user created event: %v", err)
		}
	}

	return createdUser, nil
}

// CheckLogin 驗證用戶登入
func (s *UserService) CheckLogin(req request.UserLoginRequest) (*string, error) {
	// 先驗證帳號格式（快速失敗）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, fmt.Errorf("invalid account format: %w", err)
	}

	// 查詢用戶
	user, err := s.repo.GetUsersByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	// 驗證密碼
	if !user.VerifyPassword(req.Password) {
		return nil, fmt.Errorf("password incorrect")
	}

	// 生成 JWT Token
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginUser, user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	return &jwtToken, nil
}

// publishUserCreatedEvent 發布用戶創建事件
func (s *UserService) publishUserCreatedEvent(ctx context.Context, user *model.User) error {
	payload := events.UserCreatedEventPayload{
		UserID:   user.ID,
		Account:  user.GetAccount(),
		CreateAt: time.Now().Format(time.RFC3339),
	}

	evt, err := event.NewEvent(events.UserCreatedEventType, payload)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	evt.Source = "user-service"

	if err := s.publisher.Publish(ctx, evt); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("[UserService] User created event published: UserID=%d, Account=%s", user.ID, user.GetAccount())
	return nil
}
