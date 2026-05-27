// Package service 定義用戶服務層
package service

import (
	"context"
	"errors"
	"time"

	"self_go_gin/common/msgid"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/domains/common/valueobj"
	"self_go_gin/domains/user/entity"

	"self_go_gin/domains/user/events"
	"self_go_gin/domains/user/repository"
	"self_go_gin/gin_application/handler"
	"self_go_gin/infra/event"
	apperror "self_go_gin/internal/apperror"
	jwtsecret "self_go_gin/util/jwt_secret"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用戶服務層
type UserService struct {
	ctx       context.Context
	repo      repository.UserRepository
	publisher event.Publisher
}

// NewUserService 創建用戶服務層
func NewUserService(ctx context.Context, repo repository.UserRepository, publisher event.Publisher) *UserService {
	return &UserService{
		ctx:       ctx,
		repo:      repo,
		publisher: publisher,
	}
}

// CreateUser 創建用戶
func (s *UserService) CreateUser(ctx context.Context, account valueobj.Account, password valueobj.Password) (*entity.User, error) {

	// 檢查帳號是否已存在
	exist, err := s.repo.GetUsersByAccount(account.Value())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.CheckAccountExistFailed,
			err,
			apperror.WithLayer("UserService CreateUser() GetUsersByAccount()"),
			apperror.WithLogData(map[string]interface{}{
				"account": account.Value(),
			}),
		)
	}
	if exist != nil {
		// 帳號已存在
		return nil, apperror.NewAppError(
			msgid.ResourceExist,
			appmsg.AccountAlreadyExists,
			handler.ErrResourceExist,
			apperror.WithLayer("UserService CreateUser() GetUsersByAccount() ResourceExist"),
			apperror.WithLogData(map[string]interface{}{
				"account": account.Value(),
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
func (s *UserService) CheckLogin(account valueobj.Account, password string) (*string, error) {
	// 查詢用戶
	user, err := s.repo.GetUsersByAccount(account.Value())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewAppError(
				msgid.LoginNotExist,
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
	if !user.VerifyPassword(password) {
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

	err = s.publishUserCheckLoginEvent(user)
	if err != nil {
		// 記錄錯誤但不阻止登入流程
		return nil, apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPublishFailed,
			err,
			apperror.WithLayer("UserService CheckLogin() publishUserCheckLoginEvent()"),
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

	evt, err := event.NewEvent(events.UserCreatedEventType, payload, "userService: publishUserCreatedEvent", ctx.Value("traceID").(string), time.Now())
	if err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPublishFailed,
			err,
			apperror.WithLayer("UserService publishUserCreatedEvent() NewEvent()"),
		)
	}

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

// publishUserCheckLoginEvent 發布用戶登入事件
func (s *UserService) publishUserCheckLoginEvent(user *entity.User) error {
	payload := events.UserCheckLoginEventPayload{
		UserID:  user.ID,
		Account: user.GetAccount(),
		LoginAt: time.Now().Format(time.RFC3339),
	}

	evt, err := event.NewEvent(events.UserCheckLoginEventType, payload, "userService:publishUserCheckLoginEvent", s.ctx.Value("traceID").(string), time.Now())
	if err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPublishFailed,
			err,
			apperror.WithLayer("UserService publishUserCheckLoginEvent() NewEvent()"),
		)
	}

	if err := s.publisher.Publish(s.ctx, evt); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			"failed to publish event",
			err,
			apperror.WithLayer("UserService publishUserCheckLoginEvent() Publish()"),
		)
	}

	zap.S().Infof("User check login event published successfully for UserID: %d, Account: %s", user.ID, user.GetAccount())
	return nil
}
