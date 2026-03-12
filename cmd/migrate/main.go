package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"self_go_gin/gin_application/router"
	validlang "self_go_gin/gin_application/validate_lang"
	"self_go_gin/infra/database/migrate"
	"self_go_gin/infra/database/seeder"
	"self_go_gin/infra/env"
	"self_go_gin/infra/orm/gorm_mysql"
	"self_go_gin/util/jwt_secret"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	serverEnv  = &env.ServerConfig{}
	withSeeder = flag.Bool("with-seeder", false, "Run seeder after migration")
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
	flag.Parse()

	fmt.Println("=================================")
	fmt.Println("Database Migration Tool")
	fmt.Println("=================================")

	initSetting()

	fmt.Println("Running database migration...")
	migrate.Migrate() // migrate database
	fmt.Println("✓ Migration completed successfully!")

	if *withSeeder {
		fmt.Println("Running database seeder...")
		seeder.RunSeeder() // create seeder data
		fmt.Println("✓ Seeder completed successfully!")
	}

	fmt.Println("=================================")
	fmt.Println("All tasks completed!")
	fmt.Println("=================================")
	// httpServerRun()

	//測試 log 切割
	// for i := 0; i < 2000; i++ {
	// 	wg.Add(2)
	// 	go simpleHttpGet("www.baidu.com")
	// 	go simpleHttpGet("https://www.baidu.com")
	// }
	// wg.Wait()
}

func httpServerRun() {
	quit := make(chan os.Signal, 1)
	// Set Router
	router := router.Router(quit)
	// Listen and Server
	// serverPort := ":" + strconv.Itoa(initialize.GetServerEnv().GetServerPort())

	//優雅的關閉服務(服務端關機命令發出後不會立即關機)
	//建立一個http.Server
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(GetServerEnv().Port),
		Handler: router,
	}

	go func() {
		//啟動 http.Server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Error("Server listen : %s\n", err)
		}
	}()

	/*
		監聽等待 SIGINT 或 SIGTERM 信号
		SIGINT -> 由使用者在終端中按下 Ctrl+C 產生，用於請求進程中斷
		SIGTERM -> 系統預設的終止信號，當你使用 kill 命令（不帶任何信號選項）
	*/
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Preparing Shutdown Server ...")
	//建立超時上下文，Shutdown可以讓未處理的連線在這個時間內關閉
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//停止HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Error("Shutdown Error :", err)
	}

	zap.S().Info("Server exiting")
}

func initSetting() {
	// 支持 Docker 環境和本地開發環境
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("Config path: %s\n", configPath)
	serverEnv := env.GetConfigManager().GetServerEnv()
	err := env.InitEnv(configPath)
	if err != nil {
		cfgFile := filepath.Join(configPath, "env.yaml")
		fmt.Fprintf(os.Stderr, "配置初始化失败: %v\n", err)
		fmt.Fprintf(os.Stderr, "期望的配置文件路径: %s\n", cfgFile)
		os.Exit(1)
	}
	fmt.Printf("配置信息 : %+v\n", serverEnv)
	gin.SetMode(serverEnv.AppMode)
	gorm_mysql.InitMysql(serverEnv)
	// Redis is optional for migration
	// redis.InitRedis(GetServerEnv)
	jwt_secret.SetJwtSecret(serverEnv.JwtSecret)
	// vaildate 中文化
	if err := validlang.InitValidateLang("zh"); err != nil {
		fmt.Fprintf(os.Stderr, "init trans failed, err:%v\n", err)
		panic(err)
	}
}

// GetServerEnv 獲取服務配置
func GetServerEnv() *env.ServerConfig {
	return serverEnv
}
