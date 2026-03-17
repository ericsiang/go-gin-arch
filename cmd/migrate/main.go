// Package main 是 migrate 服务的入口，建立資料表跟初始資料
package main

import (
	"flag"
	"fmt"
	"os"

	"self_go_gin/container"
	"self_go_gin/infra/database/migrate"
	"self_go_gin/infra/database/seeder"
)

var (
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

	// 1. 获取配置路径
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("Config path: %s\n", configPath)

	// 2. 初始化容器（只需要数据库连接）
	app := container.GetContainer()
	if err := app.Initialize(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Running database migration...")
	migrate.Migrate() // migrate database
	fmt.Println("✓ Migration completed successfully!")

	if *withSeeder {
		fmt.Println("Running database seeder...")
		seeder.RunSeeder() // create seeder data
		fmt.Println("✓ Seeder completed successfully!")
	}

	// 4. 清理资源
	if err := app.Shutdown(); err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
	}

	fmt.Println("=================================")
	fmt.Println("All tasks completed!")
	fmt.Println("=================================")
}
