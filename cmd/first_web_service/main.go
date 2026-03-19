// Package main 是 web 服务的入口，负责初始化配置、数据库连接、路由设置，并启动 HTTP 服务器。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	constset "self_go_gin/common/const"
	"self_go_gin/container"

	"self_go_gin/gin_application/router"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title  Self go gin Swagger API
// @version 1.0
// @description swagger first example
// @host localhost:5000
// @accept 		json
// @schemes 	http https
// @securityDefinitions.apikey	JwtTokenAuth
// @in			header
// @name   		Authorization
// @description Use Bearer JWT Token
func main() {
	// 1. 获取配置路径
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("Config path: %s\n", configPath)

	// 2. 初始化容器（统一管理所有依赖）
	app := container.GetContainer()
	if err := app.Initialize(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// 3. 初始化通用组件（JWT、验证器等）
	if err := container.InitCommonComponents(app.GetConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize common components: %v\n", err)
		os.Exit(1)
	}

	// 4. 运行服务器
	httpServerRun(app)
}

func httpServerRun(app *container.AppContainer) {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 设置路由
	router := router.Router()
	config := app.GetConfig()

	addr := ":" + strconv.Itoa(config.Port)
	// 创建 HTTP Server（添加超时配置）
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second, // 防止慢速客户端攻击
		WriteTimeout: 15 * time.Second, // 防止慢速响应
		IdleTimeout:  60 * time.Second, // Keep-Alive 超时
	}

	// 在 goroutine 中启动服务器
	go func() {
		fmt.Printf("HTTP Server is ready and listening on %s\n", addr)
		fmt.Printf("Swagger UI: http://localhost:%d/swagger-test/index.html\n", config.Port)
		// 服務連線
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 服务器启动失败
			fmt.Fprintf(os.Stderr, "Server Error: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号或服务器错误
	sig := <-quit
	// 接收到关闭信号
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Received signal: %s\n", sig)
	fmt.Println("Initiating graceful shutdown...")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// 優雅關閉流程
	gracefulShutdown(srv, app)

	fmt.Println("\n Server exited gracefully")
}

// gracefulShutdown 執行優雅關閉流程
func gracefulShutdown(srv *http.Server, app *container.AppContainer) {
	fmt.Println("Starting graceful shutdown...")
	// 創建關閉上下文（設定超時時間）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), constset.ShutdownTimeout)
	defer shutdownCancel()

	// 停止接受新 request ，並在 shutdownCtx 超時時間內等待已接收的request 請求完成
	fmt.Println("Stopping accepting new connections and waiting for active requests to complete...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Server shutdown error: %v\n", err)
		zap.L().Error("HTTP server shutdown error", zap.Error(err))
	} else {
		fmt.Println("HTTP Server shutdown completed (all requests drained)")
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
		// 在某些環境（如 stdout/stderr）中 Sync 可能會返回錯誤，這是正常的
		fmt.Printf("Log sync warning: %v (this is usually harmless)\n", err)
	} else {
		fmt.Println("Logs flushed")
	}
}
