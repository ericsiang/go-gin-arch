package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	validlang "self_go_gin/gin_application/validate_lang"
	"self_go_gin/infra/database/migrate"
	"self_go_gin/infra/database/seeder"
	"self_go_gin/infra/env"
	"self_go_gin/infra/orm/gorm_mysql"
	"self_go_gin/util/jwt_secret"

	"github.com/gin-gonic/gin"
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
