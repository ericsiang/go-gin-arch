// Package service 定義用戶服務層
package service

import (
	"context"
	"errors"
	"time"

	"self_go_gin/common/msgid"
	"self_go_gin/container"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/domains/user/entity"

	"self_go_gin/domains/user/events"
	"self_go_gin/domains/user/repository"
	"self_go_gin/gin_application/api/v1/user/request"
	"self_go_gin/gin_application/handler"
	"self_go_gin/infra/event"
	apperror "self_go_gin/internal/apperror"
	jwtsecret "self_go_gin/util/jwt_secret"

	"go.uber.org/zap"
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
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.InitFailed,
			err,
			apperror.WithLayer("UserService NewUserService()"),
		)
	}
	app := container.GetContainer()
	if app.GetConfig().IsEventBroker {
		broker := app.GetEventBroker()
		return &UserService{
			repo:      repo,
			publisher: broker.Publisher(), // 使用工廠獲取事件發布器
		}, nil
	}

	return &UserService{
		repo: repo,
	}, nil
}

// CreateUser 創建用戶
func (s *UserService) CreateUser(ctx context.Context, req request.CreateUserRequest) (*entity.User, error) {
	// 創建帳號值物件（自動驗證格式）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AccountFormatInvalid,
			err,
			apperror.WithLayer("UserService CreateUser() NewAccount()"),
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
			apperror.WithLayer("UserService CreateUser() NewPasswordFromPlainText()"),
		)
	}


	// 檢查帳號是否已存在
	_, err = s.repo.GetUsersByAccount(req.Account)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CheckAccountExistFailed,
			err,
			apperror.WithLayer("UserService CreateUser() GetUsersByAccount()"),
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
			apperror.WithLayer("UserService CreateUser() GetUsersByAccount() ResourceExist"),
			apperror.WithLogData(map[string]interface{}{
				"account": req.Account,
			}),
		)
	}

	// 創建聚合根
	user := entity.NewUser(account, password)

	// 儲存到資料庫
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CreateFailed,
			err,
			apperror.WithLayer("UserService CreateUser() CreateUser()"),
			apperror.WithLogData(map[string]interface{}{
				"user": user,
			}),
		)
	}

	// 發布用戶創建事件
	if s.publisher != nil {
		if err := s.publishUserCreatedEvent(ctx, createdUser); err != nil {
			// 記錄錯誤但不阻止用戶創建流程
			zap.S().Error("Failed to publish user created event ", zap.Error(err))
		}
	}

	return createdUser, nil
}

// CheckLogin 驗證用戶登入
func (s *UserService) CheckLogin(req request.UserLoginRequest) (*string, error) {
	// 先驗證帳號格式（快速失敗）
	account, err := valueobj.NewAccount(req.Account)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.AccountFormatInvalid,
			err,
			apperror.WithLayer("UserService CheckLogin() NewAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": req.Account,
			}),
		)
	}

	// 查詢用戶
	user, err := s.repo.GetUsersByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewAppError(
				msgid.NoContent,
				appmsg.RecordNotFound,
				handler.ErrRecordNotFound,
				apperror.WithLayer("UserService CheckLogin() GetUsersByAccount() RecordNotFound"),
				apperror.WithLogData(map[string]interface{}{
					"account": account.Value(),
				}),
			)
		}
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.QueryFailed,
			err,
			apperror.WithLayer("UserService CheckLogin() GetUsersByAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account.Value(),
			}),
		)
	}

	// 驗證密碼
	if !user.VerifyPassword(req.Password) {
		return nil, apperror.NewAppError(
			msgid.InvalidInput,
			appmsg.PasswordIncorrect,
			handler.ErrPasswordIncorrect,
			apperror.WithLayer("UserService CheckLogin() VerifyPassword()"),
		)
	}

	// 生成 JWT Token
	jwtToken, err := jwtsecret.GenerateToken(jwtsecret.LoginUser, user.ID)
	if err != nil {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.TokenGenerateFailed,
			err,
			apperror.WithLayer("UserService CheckLogin() GenerateToken()"),
		)
	}

	return &jwtToken, nil
}

// publishUserCreatedEvent 發布用戶創建事件
func (s *UserService) publishUserCreatedEvent(ctx context.Context, user *entity.User) error {
	payload := events.UserCreatedEventPayload{
		UserID:   user.ID,
		Account:  user.GetAccount(),
		CreateAt: time.Now().Format(time.RFC3339),
	}

	evt, err := event.NewEvent(events.UserCreatedEventType, payload)
	if err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPublishFailed,
			err,
			apperror.WithLayer("UserService publishUserCreatedEvent() NewEvent()"),
		)
	}

	evt.Source = "user-service"

	if err := s.publisher.Publish(ctx, evt); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			"failed to publish event",
			err,
			apperror.WithLayer("UserService publishUserCreatedEvent() Publish()"),
		)
	}

	zap.S().Infof("User created event published successfully for UserID: %d, Account: %s", user.ID, user.GetAccount())
	return nil
}
