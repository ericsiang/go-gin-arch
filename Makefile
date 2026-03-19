.PHONY: help build up down logs clean restart ps migrate migrate-seed

help: ## 顯示幫助信息
	@echo "可用的命令："
	@echo "  make build        		- 構建 Docker 映像"
	@echo "  make up          		- 啟動所有服務"
	@echo "  make down         		- 停止所有服務"
	@echo "  make logs         		- 查看日誌"
	@echo "  make clean        		- 清理所有容器和卷"
	@echo "  make restart      		- 重啟服務"
	@echo "  make ps           		- 查看運行中的容器"
	@echo "  make migrate      		- 容器執行資料庫遷移（初始化表結構）"
	@echo "  make migrate-seed	 	- 容器執行資料庫遷移並填充種子資料"
	@echo "  make run-web      		- web 本地運行"
	@echo "  make run-event-worker 	- event worker 本地運行"
	@echo "  make run-migrate  		- 資料庫遷移本地運行"



build: ## 構建 Docker 映像
	@echo "構建 Docker 映像..."
	cd scripts/docker && docker-compose build

up: ## 啟動所有服務
	@echo "啟動所有服務..."
	cd scripts/docker && docker-compose up -d
	@echo "服務已啟動！"
down: ## 停止所有服務
	@echo "停止所有服務..."
	cd scripts/docker && docker-compose down

logs: ## 查看日誌
	cd scripts/docker && docker-compose logs -f app

logs-all: ## 查看所有服務日誌
	cd scripts/docker && docker-compose logs -f

clean: ## 清理所有容器和卷
	@echo "清理所有容器和卷..."
	cd scripts/docker && docker-compose down -v
	@echo "清理完成！"

restart: ## 重啟服務
	@echo "重啟服務..."
	cd scripts/docker && docker-compose restart app

ps: ## 查看運行中的容器
	cd scripts/docker && docker-compose ps

# 資料庫遷移
migrate: ## 執行資料庫遷移（初始化資料表）
	@echo "執行資料庫遷移..."
	@docker exec gin-app /app/migrate
	@echo "✓ 資料庫遷移完成！"

migrate-seed: ## 執行資料庫遷移並填充種子資料
	@echo "執行資料庫遷移並填充種子資料..."
	@docker exec gin-app /app/migrate --with-seeder
	@echo "✓ 資料庫遷移和種子資料填充完成！"

migrate-local: ## 本地執行資料庫遷移
	@echo "本地執行資料庫遷移..."
	@go run cmd/migrate/main.go
	@echo "✓ 本地資料庫遷移完成！"

# 本地開發運行
run-web: ## 本地運行應用（不使用 Docker）
	@echo "本地運行應用..."
	go run cmd/first_web_service/main.go
run-event-worker: ## 本地運行事件工作器
	@echo "本地運行事件工作器..."
	go run cmd/event_worker/main.go
run-migrate: ## 本地運行資料庫遷移
	@echo "本地運行資料庫遷移..."
	go run cmd/migrate/main.go

# 構建本地二進制文件
build-local: ## 構建本地二進制文件
	@echo "構建本地二進制文件..."
	go build -o bin/first_web_service cmd/first_web_service/main.go
	@echo "構建完成: bin/first_web_service"

# macOS 專用二進制構建
build-mac: ## 構建 macOS 二進制文件
	@echo "構建 macOS 二進制文件..."
	GOOS=darwin GOARCH=amd64 go build -o bin/first_web_service-mac-amd64 cmd/first_web_service/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/first_web_service-mac-arm64 cmd/first_web_service/main.go
	@echo "構建完成:"
	@echo "  - Intel: bin/first_web_service-mac-amd64"
	@echo "  - Apple Silicon: bin/first_web_service-mac-arm64"
