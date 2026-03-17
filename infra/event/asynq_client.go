// Package event 提供基於 Asynq 的事件發布客戶端實現
package event

import (
	"context"
	"encoding/json"
	"fmt"
	"self_go_gin/infra/env"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
)

const (
	// DefaultQueue 默認隊列名稱
	DefaultQueue = "default"
	// HighPriorityQueue 高優先級隊列
	HighPriorityQueue = "high"
	// LowPriorityQueue 低優先級隊列
	LowPriorityQueue = "low"
)

// AsynqClient Asynq 客戶端，實現 EventPublisher 接口
type AsynqClient struct {
	client *asynq.Client
}

var asynqClient *AsynqClient

// InitAsynqClient 初始化 Asynq 客戶端
func InitAsynqClient(serverConfig *env.ServerConfig) *AsynqClient {
	redisConfig := serverConfig.Redis
	redisAddr := redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port)

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisConfig.Password,
		DB:       0,
	})

	asynqClient = &AsynqClient{
		client: client,
	}

	fmt.Println("Asynq client initialized successfully")
	return asynqClient
}

// GetAsynqClient 獲取 Asynq 客戶端實例
func GetAsynqClient() *AsynqClient {
	return asynqClient
}

// Publish 發布事件到隊列
func (c *AsynqClient) Publish(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:    DefaultQueue,
		MaxRetry: 3,
		Priority: 5,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishWithOptions 使用自定義選項發布事件（實現 EventPublisher 接口）
func (c *AsynqClient) PublishWithOptions(ctx context.Context, event *Event, opts *PublishOptions) error {
	if opts == nil {
		opts = &PublishOptions{
			Queue:    DefaultQueue,
			MaxRetry: 3,
			Priority: 5,
		}
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 構建 Asynq 選項
	asynqOpts := []asynq.Option{
		asynq.Queue(opts.Queue),
		asynq.MaxRetry(opts.MaxRetry),
	}

	// 添加延遲
	if opts.Delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(opts.Delay))
	}

	// 添加超時
	if opts.Timeout > 0 {
		asynqOpts = append(asynqOpts, asynq.Timeout(opts.Timeout))
	}

	task := asynq.NewTask(event.Type, payload, asynqOpts...)

	info, err := c.client.EnqueueContext(ctx, task, asynqOpts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	fmt.Printf("Event published: type=%s, queue=%s, id=%s\n", event.Type, info.Queue, info.ID)
	return nil
}

// PublishDeferred 延遲發布事件
func (c *AsynqClient) PublishDeferred(ctx context.Context, event *Event, delaySeconds int) error {
	opts := &PublishOptions{
		Queue:    DefaultQueue,
		MaxRetry: 3,
		Delay:    time.Duration(delaySeconds) * time.Second,
		Priority: 5,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishToHighPriorityQueue 發布到高優先級隊列
func (c *AsynqClient) PublishToHighPriorityQueue(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:    HighPriorityQueue,
		MaxRetry: 5,
		Priority: 10,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishToLowPriorityQueue 發布到低優先級隊列
func (c *AsynqClient) PublishToLowPriorityQueue(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:    LowPriorityQueue,
		MaxRetry: 1,
		Priority: 1,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// Close 關閉 Asynq 客戶端
func (c *AsynqClient) Close() error {
	return c.client.Close()
}
