// Package main 提供獨立的事件處理器服務入口
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	constset "self_go_gin/common/const"
	"self_go_gin/container"
	"self_go_gin/domains/user/events"
	"self_go_gin/infra/event"

	"go.uber.org/zap"
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

	// 5. 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 6. 在 goroutine 中启动服务器
	go func() {
		if err := subscriber.Run(); err != nil {
			log.Printf("Event server stopped with error: %v", err)
		}
	}()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Event worker server is running...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 7. 等待关闭信号
	sig := <-sigChan
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Received signal: %s\n", sig)
	fmt.Println("Initiating graceful shutdown...")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// 執行優雅關閉
	gracefulShutdownWorker(subscriber, app)

	fmt.Println("\nEvent worker server stopped gracefully")
}

// gracefulShutdownWorker 執行 event worker 的優雅關閉流程
func gracefulShutdownWorker(subscriber event.Subscriber, app *container.AppContainer) {
	fmt.Println("Graceful Shutdown ...")
	// 創建關閉上下文（總超時時間 30 秒）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), constset.ShutdownTimeout)
	defer shutdownCancel()

	fmt.Println("Stopping event subscriber...")
	if err := subscriber.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Subscriber shutdown warning: %v\n", err)
		zap.L().Warn("Event subscriber shutdown warning", zap.Error(err))
	} else {
		fmt.Println("Event subscriber stopped")
	}

	fmt.Println("Cleaning up application resources...")
	if err := app.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Resource cleanup error: %v\n", err)
		zap.L().Error("Resource cleanup error", zap.Error(err))
	} else {
		fmt.Println("All resources cleaned up")
	}

	fmt.Println("Flushing logs...")
	if err := zap.L().Sync(); err != nil {
		fmt.Printf("Log sync warning: %v (this is usually harmless)\n", err)
	} else {
		fmt.Println("Logs flushed")
	}
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
