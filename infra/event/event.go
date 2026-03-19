// Package event 提供事件驅動架構的核心定義和接口
package event

import (
	"context"
	"encoding/json"
	"time"
)

// Event 表示系統中的一個事件
type Event struct {
	// ID 事件唯一標識（可選）
	ID string `json:"id,omitempty"`
	// Type 事件類型
	Type string `json:"type"`
	// Payload 事件負載數據
	Payload json.RawMessage `json:"payload"`
	// Source 事件來源
	Source string `json:"source,omitempty"`
	// TraceID 追蹤 ID，用於分布式追蹤
	TraceID string `json:"trace_id,omitempty"`
	// Timestamp 事件時間戳
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Metadata 元數據
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Handler 事件處理器接口
type Handler interface {
	// Handle 處理事件
	Handle(ctx context.Context, event *Event) error
	// EventType 返回該處理器處理的事件類型
	EventType() string
}

// PublishOptions 發布選項
type PublishOptions struct {
	// Queue 隊列名稱
	Queue string
	// Priority 優先級 (1-10, 10 最高)
	Priority int
	// MaxRetry 最大重試次數
	MaxRetry int
	// Delay 延遲發布時間
	Delay time.Duration
	// Timeout 處理超時時間
	Timeout time.Duration
}

// Publisher 事件發布者接口
type Publisher interface {
	// Publish 發布事件
	Publish(ctx context.Context, event *Event) error
	// PublishWithOptions 使用自定義選項發布事件
	PublishWithOptions(ctx context.Context, event *Event, opts *PublishOptions) error
	// PublishDeferred 延遲發布事件（簡化方法）
	PublishDeferred(ctx context.Context, event *Event, delaySeconds int) error
	// Close 關閉發布者
	Close() error
}

// Subscriber 事件訂閱者接口
type Subscriber interface {
	// Subscribe 訂閱事件
	Subscribe(handler Handler) error
	// Start 啟動訂閱者（非阻塞）
	Start() error
	// Run 運行訂閱者（阻塞）
	Run() error
	// Shutdown 優雅關閉訂閱者
	Shutdown(ctx context.Context) error
}


// NewEvent 創建一個新事件
func NewEvent(eventType string, payload interface{}) (*Event, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type:    eventType,
		Payload: payloadBytes,
	}, nil
}

// UnmarshalPayload 解析事件負載到指定的結構體
func (e *Event) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}
