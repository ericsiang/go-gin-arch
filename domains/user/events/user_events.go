// Package events 定義 User 領域相關的事件
package events

import (
	"context"
	"fmt"
	"log"

	"self_go_gin/infra/event"
)

// 定義事件類型常量
const (
	// UserCreatedEventType 用戶創建事件
	UserCreatedEventType = "user.created"
	// UserUpdatedEventType 用戶更新事件
	UserUpdatedEventType = "user.updated"
	// UserDeletedEventType 用戶刪除事件
	UserDeletedEventType = "user.deleted"
)

// UserCreatedEventPayload 用戶創建事件的負載
type UserCreatedEventPayload struct {
	UserID   uint   `json:"user_id"`
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
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// 這裡實現具體的業務邏輯
	// 例如：發送歡迎郵件、記錄審計日誌、更新統計數據等
	log.Printf("[UserCreatedEvent] Processing user creation - UserID: %d, Account: %s",
		payload.UserID, payload.Account)

	log.Printf("[UserCreatedEvent] Successfully processed user creation event for UserID: %d", payload.UserID)
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
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[UserUpdatedEvent] Processing user update - UserID: %d, Account: %s",
		payload.UserID, payload.Account)
	// 實現具體的業務邏輯

	log.Printf("[UserUpdatedEvent] Successfully processed user update event for UserID: %d", payload.UserID)
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
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[UserDeletedEvent] Processing user deletion - UserID: %d, Account: %s",
		payload.UserID, payload.Account)
	// 實現具體的業務邏輯

	log.Printf("[UserDeletedEvent] Successfully processed user deletion event for UserID: %d", payload.UserID)
	return nil
}

