# Docker 部署說明

本專案已配置完整的 Docker 容器化部署方案，支持在 macOS 上快速啟動。

## 📁 Docker 文件結構

```
scripts/docker/
├── Dockerfile              # Docker 映像構建文件
├── docker-compose.yaml     # Docker Compose 配置文件
└── DOCKER_DEPLOYMENT.md    # 本說明文件
```

## 🚀 快速開始

### 前置需求
- Docker Desktop for Mac
- Docker Compose (通常已包含在 Docker Desktop 中)

### 使用 Makefile 命令（推薦）

在專案根目錄執行：

```bash
# 1. 構建 Docker 映像
make build

# 2. 啟動所有服務（包括 MySQL, Redis, App）
make up

# 3. 執行資料庫遷移（首次啟動必須執行）
make migrate

# 4. 查看應用日誌
make logs

# 5. 停止所有服務
make down
```

**注意**: 首次啟動後必須執行 `make migrate` 來初始化資料表結構。

### 使用 Docker Compose 命令

```bash
# 進入 docker 目錄
cd scripts/docker

# 構建並啟動所有服務
docker-compose up -d

# 查看容器狀態
docker-compose ps

# 查看日誌
docker-compose logs -f app

# 停止服務
docker-compose down
```

## 📦 容器架構

專案包含三個服務：

1. **gin-app** - Go Gin 應用程式
   - Port: 5001
   - 基於 Alpine Linux
   - 多階段構建優化鏡像大小

2. **gin-mysql** - MySQL 數據庫
   - Port: 3307（主機訪問）
   - 內部 Port: 3306（容器間通信）
   - 用戶: root
   - 密碼: 123456
   - 數據庫: siang_gin

3. **gin-redis** - Redis 緩存
   - Port: 6378（主機訪問）
   - 內部 Port: 6379（容器間通信）
   - 無密碼

## 🎯 完整初始化流程

首次部署時，請按照以下步驟操作：

```bash
# 1. 構建 Docker 映像
make build

# 2. 啟動所有服務
make up

# 3. 等待服務啟動（約 10-15 秒）
sleep 15

# 4. 執行資料庫遷移（建立資料表）
make migrate

# 5. 查看應用日誌確認正常運行
make logs

# 6. 訪問應用
# 瀏覽器打開: http://localhost:5001
```

如果需要測試資料，可以執行：
```bash
make migrate-seed
```

## 🔧 配置文件

- **本地開發**: `conf/env.yaml`
- **Docker 環境**: `conf/env.docker.yaml`

Docker 容器會自動使用 `env.docker.yaml`，其中數據庫和 Redis 的 host 設為對應的服務名稱。

## 📝 常用命令

在專案根目錄執行：

```bash
# 查看所有可用命令
make help

# 重啟應用
make restart

# 查看所有服務日誌
make logs-all

# 完全清理（包括數據卷）
make clean

# 檢查運行狀態
make ps
```

## 🗄️ 資料庫遷移

首次啟動容器後，需要執行資料庫遷移來初始化資料表結構。

### 初始化資料表

```bash
# 僅執行資料庫遷移（建立資料表）
make migrate

# 或手動執行
docker exec gin-app /app/migrate
```

### 遷移並填充測試資料

如果需要同時填充種子資料（測試資料）：

```bash
# 執行遷移並填充種子資料
make migrate-seed

# 或手動執行
docker exec gin-app /app/migrate --with-seeder
```

### 本地環境遷移

如果在本地環境（非 Docker）運行：

```bash
# 本地執行遷移
make migrate-local

# 或直接運行
go run cmd/migrate/main.go
```

### 查看遷移狀態

```bash
# 連接到 MySQL 查看資料表
docker exec -it gin-mysql mysql -uroot -p123456 siang_gin -e "SHOW TABLES;"
```

## 📝 常用 Docker 命令
make ps
```

## 🏗️ 本地構建二進制文件

如果需要在 macOS 上直接運行（不使用 Docker）：

```bash
# 構建 macOS 二進制文件（Intel 和 Apple Silicon）
make build-mac

# 或者構建當前架構版本
make build-local

# 本地運行（需要先啟動 MySQL 和 Redis）
make run-local
```

## 🔍 訪問應用

服務啟動後，可以通過以下地址訪問：

- **應用程式**: http://localhost:5001
- **Swagger 文檔** (如有): http://localhost:5001/swagger/index.html
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

## 🐛 故障排查

### 端口被占用
```bash
# 檢查端口占用
lsof -i :5001
lsof -i :3307
lsof -i :6378

# 修改 scripts/docker/docker-compose.yaml 中的端口映射
```

### 查看容器日誌
```bash
# 查看應用日誌
docker logs gin-app

# 查看 MySQL 日誌
docker logs gin-mysql

# 查看 Redis 日誌
docker logs gin-redis
```

### 進入容器內部調試
```bash
# 進入應用容器
docker exec -it gin-app sh

# 進入 MySQL 容器
docker exec -it gin-mysql mysql -uroot -p123456

# 進入 Redis 容器
docker exec -it gin-redis redis-cli
```

### 資料庫相關問題

**問題: 應用啟動失敗，提示資料表不存在**
```bash
# 解決方法：執行資料庫遷移
make migrate

# 或手動執行
docker exec gin-app /app/migrate
```

**問題: 資料庫連接失敗**
```bash
# 檢查 MySQL 容器是否正常運行
docker ps | grep mysql

# 檢查 MySQL 健康狀態
docker inspect gin-mysql | grep -A 10 Health

# 查看 MySQL 日誌
docker logs gin-mysql --tail 50

# 手動測試連接
docker exec -it gin-mysql mysql -uroot -p123456 -e "SELECT 1;"
```

**問題: 需要重建資料表**
```bash
# 方法 1: 刪除 MySQL volume 並重新初始化
make down
docker volume rm docker_mysql_data
make up
sleep 15
make migrate

# 方法 2: 手動刪除資料表
docker exec -it gin-mysql mysql -uroot -p123456 siang_gin -e "DROP TABLE IF EXISTS users, admins;"
make migrate
```

**問題: 查看已建立的資料表**
```bash
# 列出所有資料表
docker exec -it gin-mysql mysql -uroot -p123456 siang_gin -e "SHOW TABLES;"

# 查看資料表結構
docker exec -it gin-mysql mysql -uroot -p123456 siang_gin -e "DESCRIBE users;"
```

## 📊 健康檢查

Docker Compose 已配置健康檢查：
- MySQL: 每 10 秒檢查一次
- Redis: 每 10 秒檢查一次
- App: 依賴 MySQL 和 Redis 健康後才啟動

## 🔄 更新代碼後重新部署

```bash
# 方法 1: 使用 Makefile（推薦）
make build
make up

# 方法 2: 使用 Docker Compose
cd scripts/docker && docker-compose up -d --build
```

## 🗑️ 清理環境

```bash
# 停止服務但保留數據
make down

# 完全清理（包括數據庫數據）
make clean
```

## 📌 注意事項

1. **首次啟動**: 必須執行 `make migrate` 來初始化資料表，否則應用會因找不到資料表而無法正常運行
2. **數據持久化**: 數據持久化在 Docker volumes 中，停止容器不會丟失數據
3. **完全清理**: 使用 `make clean` 會刪除所有數據，請謹慎使用
4. **日誌文件**: 會保存在 `app_logs` volume 中
5. **Docker 文件位置**: 所有 Docker 相關文件統一放在 `scripts/docker/` 目錄下
6. **端口說明**: 
   - 主機訪問：MySQL 使用 3307，Redis 使用 6378
   - 容器間通信：MySQL 使用 3306，Redis 使用 6379

## 🆘 需要幫助？

如遇問題，請檢查：
1. Docker Desktop 是否正常運行
2. 端口是否被占用
3. 查看容器日誌獲取錯誤信息
4. 確認在正確的目錄執行命令
