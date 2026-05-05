// Package events 定義 User 領域相關的事件
package events

import (
	"context"
	"fmt"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/common/appmsg"
	"self_go_gin/infra/event"
	apperror "self_go_gin/internal/apperror"

	"go.uber.org/zap"
)

// 定義事件類型常量
const (
	// UserCreatedEventType 用戶創建事件
	UserCreatedEventType = "user.created"
	// UserUpdatedEventType 用戶更新事件
	UserUpdatedEventType = "user.updated"
	// UserDeletedEventType 用戶刪除事件
	UserDeletedEventType = "user.deleted"
	// UserCheckLoginEventType 用戶登錄事件
	UserCheckLoginEventType = "user.check_login"
)

// UserCreatedEventPayload 用戶創建事件的負載
type UserCreatedEventPayload struct {
	UserID   uint64 `json:"user_id"`
	Account  string `json:"account"`
	Email    string `json:"email,omitempty"`
	CreateAt string `json:"create_at"`
}

// UserUpdatedEventPayload 用戶更新事件的負載
type UserUpdatedEventPayload struct {
	UserID   uint   `json:"user_id"`
	Account  string `json:"account"`
	UpdateAt string `json:"update_at"`
}

// UserDeletedEventPayload 用戶刪除事件的負載
type UserDeletedEventPayload struct {
	UserID   uint   `json:"user_id"`
	Account  string `json:"account"`
	DeleteAt string `json:"delete_at"`
}

// UserCheckLoginEventPayload 用戶登錄事件的負載
type UserCheckLoginEventPayload struct {
	UserID  uint64 `json:"user_id"`
	Account string `json:"account"`
	LoginAt string `json:"login_at"`
}

// UserCreatedEventHandler 處理用戶創建事件
type UserCreatedEventHandler struct{}

// NewUserCreatedEventHandler 創建用戶創建事件處理器
func NewUserCreatedEventHandler() *UserCreatedEventHandler {
	return &UserCreatedEventHandler{}
}

// EventType 返回處理的事件類型
func (h *UserCreatedEventHandler) EventType() string {
	return UserCreatedEventType
}

// Handle 處理用戶創建事件
func (h *UserCreatedEventHandler) Handle(_ context.Context, evt *event.Event) error {
	var payload UserCreatedEventPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPayloadParseFailed,
			fmt.Errorf("failed to unmarshal payload: %w", err),
			apperror.WithLayer("Event"),
		)
	}

	// 這裡實現具體的業務邏輯
	// 例如：發送歡迎郵件、記錄審計日誌、更新統計數據等
	zap.S().Infof("Processing user creation event for UserID: %d, Account: %s", payload.UserID, payload.Account)

	zap.S().Infof("User created event processed successfully for UserID: %d", payload.UserID)
	fmt.Printf("Processing user creation event: payload=%+v:\n", payload)

	return nil
}

// UserUpdatedEventHandler 處理用戶更新事件
type UserUpdatedEventHandler struct{}

// NewUserUpdatedEventHandler 創建用戶更新事件處理器
func NewUserUpdatedEventHandler() *UserUpdatedEventHandler {
	return &UserUpdatedEventHandler{}
}

// EventType 返回處理的事件類型
func (h *UserUpdatedEventHandler) EventType() string {
	return UserUpdatedEventType
}

// Handle 處理用戶更新事件
func (h *UserUpdatedEventHandler) Handle(_ context.Context, evt *event.Event) error {
	var payload UserUpdatedEventPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPayloadParseFailed,
			fmt.Errorf("failed to unmarshal payload: %w", err),
			apperror.WithLayer("Event"),
		)
	}

	zap.S().Infof("Processing user update event for UserID: %d, Account: %s", payload.UserID, payload.Account)
	// 實現具體的業務邏輯
	zap.S().Infof("User updated event processed successfully for UserID: %d", payload.UserID)
	return nil
}

// UserDeletedEventHandler 處理用戶刪除事件
type UserDeletedEventHandler struct{}

// NewUserDeletedEventHandler 創建用戶刪除事件處理器
func NewUserDeletedEventHandler() *UserDeletedEventHandler {
	return &UserDeletedEventHandler{}
}

// EventType 返回處理的事件類型
func (h *UserDeletedEventHandler) EventType() string {
	return UserDeletedEventType
}

// Handle 處理用戶刪除事件
func (h *UserDeletedEventHandler) Handle(_ context.Context, evt *event.Event) error {
	var payload UserDeletedEventPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPayloadParseFailed,
			fmt.Errorf("failed to unmarshal payload: %w", err),
			apperror.WithLayer("Event"),
		)
	}

	zap.S().Infof("Processing user deletion event for UserID: %d, Account: %s", payload.UserID, payload.Account)

	// 實現具體的業務邏輯
	zap.S().Infof("User deleted event processed successfully for UserID: %d", payload.UserID)
	return nil
}

// UserCheckLoginEventHandler 處理用戶登錄事件
type UserCheckLoginEventHandler struct{}

func NewUserCheckLoginEventHandler() *UserCheckLoginEventHandler {
	return &UserCheckLoginEventHandler{}
}

func (h *UserCheckLoginEventHandler) EventType() string {
	return UserCheckLoginEventType
}

func (h *UserCheckLoginEventHandler) Handle(_ context.Context, evt *event.Event) error {
	var payload UserCheckLoginEventPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		return apperror.NewAppError(
			msgid.Fail,
			appmsg.UserEventPayloadParseFailed,
			fmt.Errorf("failed to unmarshal payload: %w", err),
			apperror.WithLayer("Event"),
		)
	}

	fmt.Printf("UserCheckLoginEvent payload: %+v\n", payload)
	return nil
}
