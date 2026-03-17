// Package main 提供獨立的事件處理器服務入口
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"self_go_gin/container"
	"self_go_gin/domains/user/events"
	"self_go_gin/infra/event"
)

func main() {
	// 1. 获取配置路径
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("Config path: %s\n", configPath)

	// 2. 初始化容器（需要 Redis 和事件代理）
	app := container.GetContainer()
	if err := app.Initialize(configPath); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// 3. 获取事件代理
	broker := app.GetEventBroker()
	if broker == nil {
		log.Fatal("Event broker is not initialized. Please enable IsEventBroker in config.")
	}
	subscriber := broker.Subscriber()

	// 4. 注册事件处理器
	if err := registerEventHandlers(subscriber); err != nil {
		log.Fatalf("Failed to register event handlers: %v", err)
	}

	// 5. 设置优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 6. 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 7. 在 goroutine 中启动服务器
	go func() {
		if err := subscriber.Run(); err != nil {
			log.Printf("Event server stopped with error: %v", err)
		}
	}()

	fmt.Println("Event worker server is running...")
	fmt.Println("Broker type: Asynq (Redis-based)")
	fmt.Println("Press Ctrl+C to stop")

	// 8. 等待关闭信号
	<-sigChan
	fmt.Println("\nReceived shutdown signal")
	subscriber.Shutdown()

	// 9. 清理容器资源
	if err := app.Shutdown(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	<-ctx.Done()
	fmt.Println("Event worker server stopped gracefully")
}

// registerEventHandlers 註冊所有事件處理器
func registerEventHandlers(subscriber event.Subscriber) error {
	// 註冊 User 領域的事件處理器
	handlers := []event.Handler{
		events.NewUserCreatedEventHandler(),
		events.NewUserUpdatedEventHandler(),
		events.NewUserDeletedEventHandler(),
	}

	for _, handler := range handlers {
		if err := subscriber.Subscribe(handler); err != nil {
			return fmt.Errorf("failed to subscribe handler %s: %w", handler.EventType(), err)
		}
	}

	fmt.Printf("Successfully registered %d event handlers\n", len(handlers))
	return nil
}
