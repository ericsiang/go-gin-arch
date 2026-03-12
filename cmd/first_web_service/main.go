package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"self_go_gin/gin_application/router"
	validlang "self_go_gin/gin_application/validate_lang"

	"self_go_gin/infra/cache/redis"
	"self_go_gin/infra/env"
	"self_go_gin/infra/orm/gorm_mysql"
	"self_go_gin/util/jwt_secret"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
	// 启动初始化（失败会 panic，这是正确的）
	initSetting()
	// 运行服务器
	httpServerRun()
}

func httpServerRun() {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 设置路由
	router := router.Router(quit)
	serverEnv := env.GetConfigManager().GetServerEnv()

	addr := ":" + strconv.Itoa(serverEnv.Port)
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
		fmt.Printf("Swagger UI: http://localhost:%d/swagger-test/index.html\n", serverEnv.Port)
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
	fmt.Printf("\n Received signal: %s\n", sig)
	fmt.Println("Initiating graceful shutdown...")

	// 优雅关闭（5秒超时）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// 关闭 HTTP 服务器
	fmt.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Server shutdown error: %v\n", err)
	} else {
		fmt.Println("HTTP Server shutdown completed")
	}

	// 清理资源
	cleanupResources()

	fmt.Println("Server exited gracefully")
}

// cleanupResources 清理所有资源（数据库连接等）
func cleanupResources() {
	fmt.Println("Cleaning up resources...")

	// 关闭数据库连接
	if db, err := gorm_mysql.GetMysqlDB(); err != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close database: %v\n", err)
			} else {
				fmt.Println("Mysql connection closed")
			}
		}
	}

	// 这里可以添加其他资源清理逻辑

}

func initSetting() {
	fmt.Println("\n Initializing application...")
	serverEnv := env.GetConfigManager().GetServerEnv()

	// 1. 获取配置路径
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("Config path: %s\n", configPath)

	// 2. 加载配置
	err := env.InitEnv(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n 配置初始化失败: %v\n", err)
		cfgFile := filepath.Join(configPath, "env.yaml")
		fmt.Fprintf(os.Stderr, "期望的配置文件路径: %s\n", cfgFile)
		fmt.Fprintln(os.Stderr, "请检查配置文件是否存在且格式正确")
		panic(err)
	}
	fmt.Println("Configuration loaded")

	// 3. 设置 Gin 模式
	gin.SetMode(serverEnv.AppMode)
	fmt.Printf("Gin mode: %s\n", serverEnv.AppMode)

	// 4. 初始化数据库
	if serverEnv.MysqlDB.Host != "" {
		gorm_mysql.InitMysql(serverEnv)
	}

	// 5. 初始化 Redis（如果需要）
	if serverEnv.Redis.Host != "" {
		redis.InitRedis(serverEnv)
	}

	// 6. 设置 JWT 密钥
	jwt_secret.SetJwtSecret(serverEnv.JwtSecret)
	fmt.Println("JWT secret configured")

	// 7. 初始化验证器中文化
	if err := validlang.InitValidateLang("zh"); err != nil {
		fmt.Fprintf(os.Stderr, "\n 验证器初始化失败: %v\n", err)
		panic(err)
	}
	fmt.Println("Validator localization initialized")

	fmt.Println("\n All components initialized successfully!")
}
