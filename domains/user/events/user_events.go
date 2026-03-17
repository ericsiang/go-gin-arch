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

	// 示例：發送歡迎郵件（實際應該調用郵件服務）
	if err := h.sendWelcomeEmail(payload); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
		// 根據業務需求決定是否返回錯誤（返回錯誤會觸發重試）
	}

	// 示例：記錄審計日誌
	h.logAudit(payload)

	log.Printf("[UserCreatedEvent] Successfully processed user creation event for UserID: %d", payload.UserID)
	return nil
}

// sendWelcomeEmail 發送歡迎郵件（示例方法）
func (h *UserCreatedEventHandler) sendWelcomeEmail(payload UserCreatedEventPayload) error {
	// TODO: 實現實際的郵件發送邏輯
	log.Printf("Sending welcome email to user: %s", payload.Account)
	return nil
}

// logAudit 記錄審計日誌（示例方法）
func (h *UserCreatedEventHandler) logAudit(payload UserCreatedEventPayload) {
	// TODO: 實現實際的審計日誌記錄
	log.Printf("Audit log: User created - ID: %d, Account: %s", payload.UserID, payload.Account)
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
	// 例如：通知相關系統、更新緩存、記錄變更歷史等
	h.invalidateCache(payload)
	h.notifyRelatedSystems(payload)

	log.Printf("[UserUpdatedEvent] Successfully processed user update event for UserID: %d", payload.UserID)
	return nil
}

// invalidateCache 使緩存失效（示例方法）
func (h *UserUpdatedEventHandler) invalidateCache(payload UserUpdatedEventPayload) {
	// TODO: 實現緩存失效邏輯
	log.Printf("Invalidating cache for user: %d", payload.UserID)
}

// notifyRelatedSystems 通知相關系統（示例方法）
func (h *UserUpdatedEventHandler) notifyRelatedSystems(payload UserUpdatedEventPayload) {
	// TODO: 實現系統通知邏輯
	log.Printf("Notifying related systems about user update: %d", payload.UserID)
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
	// 例如：清理相關數據、通知其他系統、記錄刪除日誌等
	h.cleanupUserData(payload)
	h.logDeletion(payload)

	log.Printf("[UserDeletedEvent] Successfully processed user deletion event for UserID: %d", payload.UserID)
	return nil
}

// cleanupUserData 清理用戶數據（示例方法）
func (h *UserDeletedEventHandler) cleanupUserData(payload UserDeletedEventPayload) {
	// TODO: 實現數據清理邏輯
	log.Printf("Cleaning up data for deleted user: %d", payload.UserID)
}

// logDeletion 記錄刪除日誌（示例方法）
func (h *UserDeletedEventHandler) logDeletion(payload UserDeletedEventPayload) {
	// TODO: 實現刪除日誌記錄
	log.Printf("Audit log: User deleted - ID: %d, Account: %s", payload.UserID, payload.Account)
}
