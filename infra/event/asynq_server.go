// Package event 提供基於 Asynq 的事件處理服務器實現
package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"self_go_gin/infra/env"

	"github.com/hibiken/asynq"
)

// AsynqServer Asynq 服務器，實現 EventSubscriber 接口
type AsynqServer struct {
	server   *asynq.Server
	mux      *asynq.ServeMux
	handlers map[string]Handler
}

var asynqServer *AsynqServer

// InitAsynqServer 初始化 Asynq 服務器
func InitAsynqServer(serverConfig *env.ServerConfig) *AsynqServer {
	redisConfig := serverConfig.Redis
	redisAddr := redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisConfig.Password,
			DB:       0,
		},
		asynq.Config{
			// 設置不同隊列的優先級
			Queues: map[string]int{
				HighPriorityQueue: 6, // 高優先級隊列，權重 6
				DefaultQueue:      3, // 默認隊列，權重 3
				LowPriorityQueue:  1, // 低優先級隊列，權重 1
			},
			// 並發處理的工作數量
			Concurrency: 10,
			// 錯誤處理器
			ErrorHandler: asynq.ErrorHandlerFunc(func(_ context.Context, task *asynq.Task, err error) {
				log.Printf("Error processing task %s: %v", task.Type(), err)
			}),
		},
	)

	asynqServer = &AsynqServer{
		server:   server,
		mux:      asynq.NewServeMux(),
		handlers: make(map[string]Handler),
	}

	fmt.Println("Asynq server initialized successfully")
	return asynqServer
}

// GetAsynqServer 獲取 Asynq 服務器實例
func GetAsynqServer() *AsynqServer {
	return asynqServer
}

// Subscribe 訂閱事件，註冊事件處理器
func (s *AsynqServer) Subscribe(handler Handler) error {
	eventType := handler.EventType()

	// 檢查是否已經註冊
	if _, exists := s.handlers[eventType]; exists {
		return fmt.Errorf("handler for event type %s already registered", eventType)
	}

	// 保存處理器
	s.handlers[eventType] = handler

	// 註冊到 ServeMux
	s.mux.HandleFunc(eventType, func(ctx context.Context, task *asynq.Task) error {
		var event Event
		if err := json.Unmarshal(task.Payload(), &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		fmt.Printf("Processing event: type=%s, source=%s\n", event.Type, event.Source)
		return handler.Handle(ctx, &event)
	})

	fmt.Printf("Event handler registered: %s\n", eventType)
	return nil
}

// Start 啟動事件處理服務器
func (s *AsynqServer) Start() error {
	fmt.Println("Starting Asynq event server...")
	if err := s.server.Start(s.mux); err != nil {
		return fmt.Errorf("failed to start asynq server: %w", err)
	}
	return nil
}

// Shutdown 優雅關閉服務器
func (s *AsynqServer) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	fmt.Println("Asynq server shutdown")
	s.server.Shutdown()
	return nil
}

// Run 運行服務器（阻塞方法）
func (s *AsynqServer) Run() error {
	fmt.Println("Running Asynq event server...")
	if err := s.server.Run(s.mux); err != nil {
		return fmt.Errorf("failed to run asynq server: %w", err)
	}
	return nil
}
