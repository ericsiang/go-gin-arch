# Golang Gin Framework 設計生產級模板

## 技術棧說明
* **框架和語言**：Go 1.24+, Gin
* **資料庫**：MySQL, Redis
* **事件驅動**：Asynq, RabbitMQ
* **可觀測性**：Prometheus + Grafana
* **容器化**：Docker, Docker Compose
* **其他**：JWT, OPA, Zap, GORM, Viper, Swagger

## 專案簡介

一個**Go Web 服務模板**，基於 DDD (領域驅動設計) 和事件驅動架構，展示**生產級應用**的最佳實踐。

## 核心特色
- **可插拔的架構設計** - 應用框架可隨時替換（Gin ↔ Echo）
- **事件驅動架構** - 支援 Asynq / RabbitMQ，可透過環境變數快速切換
- **完整的 DDD 實踐** - Entity、ValueObject、Repository、Service 分層清晰
- **可觀測性** - Prometheus 指標收集 + Grafana 即時監控面板
- **生產級特性** - 統一錯誤處理、結構化日誌、Trace ID 追蹤、優雅關閉
- **高可維護性** - 業務邏輯與框架解耦，降低替換成本

## 適用場景
適合需要**長期維護**、**頻繁迭代**、**團隊協作**的中型 Web 服務項目。


## 架構設計

### 分層架構
採用 DDD 分層設計，各層職責清晰：

| 層級            | 職責         | 特點                             |
| --------------- | ------------ | -------------------------------- |
| **Entity**      | 實體與聚合根 | 封裝業務規則                     |
| **ValueObject** | 值物件       | 不可變對象，如 Account、Password |
| **Service**     | 業務邏輯協調 | 事務邊界、流程編排               |
| **Repository**  | 數據訪問抽象 | Domain ↔ Model 轉換              |
| **DAO**         | 純數據庫操作 | GORM 操作，無業務邏輯            |
| **Model**       | 數據表結構   | 持久化對象                       |

**[查看完整目錄結構](./asset/markdown/tree.md)**

### 容器層設計
- **依賴注入容器** - 統一管理 DB、Redis、EventBroker 等組件
- **生命周期管理** - Initialize() 初始化 / Shutdown() 優雅關閉
- **線程安全** - sync.RWMutex 保證並發安全
- **單例模式** - 全局唯一實例
- **可擴展** - 輕鬆添加新的基礎設施組件

## 技術特性

**基礎設施**
- 環境配置：Viper
- 日誌系統：Zap（分級日誌）
- 緩存：Redis
- 數據庫：MySQL + GORM
- 事件驅動：Asynq / RabbitMQ（可切換）

**安全機制**
- JWT 身份認證
- OPA 權限控制
- Bcrypt 密碼加密

**Web 層功能**
- 中間件：Trace ID、限流、JWT 認證、權限驗證
- 統一 Response 處理
- API 版本控制
- Swagger 文檔

**可觀測性**
- Trace 追蹤：Trace-ID 貫穿全鏈路
- 結構化日誌：Zap Logger 分級記錄
- 性能指標：Prometheus 收集 HTTP/系統指標
- 可視化監控：Grafana 儀表板展示
- 優雅關閉：確保請求處理完成

**容器化**
- Docker + Docker Compose
- 一鍵啟動開發環境

## 架構圖結構
### HTTP 請求處理流程 (DDD 分層架構)
``` mermaid
graph LR
    Client((用戶端)) -->|HTTP Request| Router["Router<br/>路由匹配"]
    
    Router --> Middleware["Middleware<br/>日誌/認證/授權"]
    
    Middleware --> Controller["Controller<br/>參數驗證與轉換"]
    
    Controller --> Service["ApplicationService<br/>事務/流程協調/業務規則檢查"]
    
    Service --> Repository["Repository<br/>領域模型 ↔ Model 轉換"]
    
    Repository --> DAO["DAO<br/>純粹數據庫操作"]
    
    DAO --> Database[("Database<br/>MySQL/Redis")]
    
    Database --> DAO
    DAO --> Repository
    Repository --> Service
    Service --> Controller
    Controller --> Response["Response<br/>統一格式輸出"]
    Response --> Client
    
```

### 框架可替換性設計
``` mermaid
graph TB
    subgraph Current ["目前架構 (使用 Gin)"]
        GinApp["gin_application"]
    end
    
    subgraph Core ["核心業務層 (框架無關)"]
        DomainCore["domains"]
        InfraCore["infra"]
    end
    
    subgraph Future ["未來可替換 (例如 Echo)"]
        EchoApp["echo_application"]
    end
    
    GinApp -.->|調用| DomainCore
    EchoApp -.->|調用| DomainCore
    DomainCore --> InfraCore
    
    Replace["🔄 替換 Web 框架<br/>只需修改 Web 框架層<br/>Domain 和 Infrastructure 層無需變動"]
    
    GinApp -.-> Replace
    Replace -.-> EchoApp
    
```

### 事件驅動架構流程
``` mermaid
graph TB
    subgraph Service ["業務服務層"]
        UserService["UserService<br/>用戶服務"]
    end
    
    subgraph EventPublish ["事件發佈"]
        Publisher["EventPublisher<br/>事件發佈者"]
        EventBroker["EventBroker<br/>事件代理器"]
    end
    
    subgraph Queue ["消息隊列 (Asynq/Redis)"]
        RedisQueue[("Redis Queue<br/>事件隊列")]
    end
    
    subgraph EventWorker ["事件處理服務 (獨立進程)"]
        Subscriber["EventSubscriber<br/>事件訂閱者"]
        Handler1["UserCreatedHandler<br/>用戶創建事件處理器"]
        Handler2["UserUpdatedHandler<br/>用戶更新事件處理器"]
        Handler3["UserDeletedHandler<br/>用戶刪除事件處理器"]
    end
    
    subgraph Actions ["異步操作"]
        SendEmail["發送郵件"]
        SyncData["同步數據"]
        UpdateCache["更新緩存"]
        Logging["記錄日誌"]
    end
    
    UserService -->|1. 業務操作完成| Publisher
    Publisher -->|2. 發佈事件| EventBroker
    EventBroker -->|3. 寫入隊列| RedisQueue
    
    RedisQueue -->|4. 消費事件| Subscriber
    Subscriber -->|5. 分發事件| Handler1
    Subscriber -->|5. 分發事件| Handler2
    Subscriber -->|5. 分發事件| Handler3
    
    Handler1 -->|6. 執行操作| SendEmail
    Handler2 -->|6. 執行操作| SyncData
    Handler3 -->|6. 執行操作| UpdateCache
    Handler3 -->|6. 執行操作| Logging
        
```

### 可觀測性架構 (Prometheus + Grafana)
``` mermaid
graph LR
    App["應用服務<br/>HTTP Handler<br/>Prometheus Middleware"]
    Prom["Prometheus<br/>:9090<br/>抓取指標"]
    Grafana["Grafana<br/>:3000<br/>視覺化監控"]
    
    subgraph Metrics["收集指標"]
        M1["HTTP 計數器<br/>請求總數、狀態"]
        M2["系統資源<br/>Goroutine、內存"]
    end
    
    App --> M1
    App --> M2

    M1 -.-> Prom
    M2 -.-> Prom

    
    Prom -->|PromQL 查詢| Grafana
    Grafana -->|即時面板| User["開發者 / PM"]
```

**核心指標**
| 指標                                      | 類型      | 場景                 |
| ----------------------------------------- | --------- | -------------------- |
| `http_requests_total`                     | Counter   | 追蹤 QPS、按端點分析 |
| `go_goroutines` / `go_memory_usage_bytes` | Gauge     | 檢測洩漏、資源告警   |

**常用 PromQL 查詢**
```promql
# 實時 QPS
rate(http_requests_total[1m])

# 實時 API status RPS
sum by (path,status) (rate(SiangGin_http_requests_total[1m]))

```

### 初始化容器層架構
``` mermaid
graph TB
    subgraph AppStart ["應用啟動流程"]
        Main["main.go"] --> Container["Container.Initialize()"]
    end
    
    subgraph ContainerLayer ["Container 容器層"]
        Container --> ConfigLoad["載入配置<br/>env.InitEnv()"]
        ConfigLoad --> InfraInit["初始化基礎設施<br/>initInfrastructure"]
        
        InfraInit --> MySQL["MySQL<br/>gormysql.InitMysql()"]
        InfraInit --> Redis["Redis<br/>redis.InitRedis()"]
        InfraInit --> EventBroker["EventBroker<br/>event.NewEventBroker()"]
        InfraInit --> Common["通用組件<br/>JWT / Validator"]
    end
    
    subgraph Usage ["各層使用容器"]
        Controller["Controller"] --> GetContainer["container.GetContainer()"]
        Service["Service"] --> GetContainer
        Repository["Repository"] --> GetContainer
        
        GetContainer --> GetDB["app.GetDB()"]
        GetContainer --> GetRedis["app.GetRedisClient()"]
        GetContainer --> GetBroker["app.GetEventBroker()"]
        GetContainer --> GetConfig["app.GetConfig()"]
    end
    
    subgraph Shutdown ["優雅關閉"]
        Signal["收到關閉信號"] --> AppShutdown["app.Shutdown()"]
        AppShutdown --> CloseRedis["關閉 redis 連接"]
        AppShutdown --> CloseBroker["關閉事件代理"]
        AppShutdown --> CloseDB["關閉數據庫連接"]
        AppShutdown --> CleanRes["清理其他資源"]
    end
    
    ContainerLayer --> Usage
    Usage --> Shutdown
    
```

## 快速開始

### 前置需求
- Go 1.24+
- Docker & Docker Compose
- golangci-lint [程式碼檢查](https://golangci-lint.run/docs/welcome/install/local/)
- wrk（壓力測試，可選）

### Docker 部署
```bash
# 1. 複製配置文件
cp conf/env.docker.yaml.example conf/env.docker.yaml
# 編輯 conf/env.docker.yaml 填入真實配置

# 2. 建立及啟動所有服務
make build && make up

# 3. 執行資料庫遷移
make migrate

# 4. 查看服務狀態
make logs
```

### 本地開發
```bash
# 1. 複製配置文件
cp conf/env.yaml.example conf/env.yaml
# 編輯 conf/env.yaml 填入真實配置

# 2. 啟動 Web 服務
make run-web

# 3. 啟動事件處理器
make run-event-worker

# 4. 資料庫遷移
make run-migrate

# 5. 生成 swagger 文檔
make swagger
```

### 可用命令
```bash
make help  # 查看所有可用命令
```

### 監控系統 (Prometheus + Grafana)
**Docker 部署包含以下服務**
- Prometheus (port 9090) - 自動抓取應用指標
- Grafana (port 3000) - 視覺化面板和告警
- 應用暴露 `/api/v1/metrics` 端點供 Prometheus 抓取

**本地開發啟動監控**
```bash
# Docker Compose 已包含 Prometheus + Grafana
make up

# 訪問:
# - Grafana: http://localhost:3000 (admin/admin)
# - Prometheus: http://localhost:9090
# - 應用指標: http://localhost:5001/api/v1/metrics
```

**監控面板內容**
- 實時 QPS
- Goroutine / 內存使用趨勢

### 壓力測試（可選）
需要安裝 wrk
```bash
# 測試建立用戶
wrk -t8 -c100 -d30s -s ./scripts/wrk/create.lua --latency http://localhost:5001

# 測試多用戶登入，先跑過測試建立用戶，接著從資料庫只搜尋 account 欄位資料，複製修改到 ./scripts/wrk/login_account.txt
wrk -t8 -c100 -d30s -s ./scripts/wrk/login.lua --latency http://localhost:5001

# 測試帶 Token 請求
wrk -t8 -c100 -d30s -s ./scripts/wrk/token.lua --latency http://localhost:5001
```

## 使用到的 package
<table>
    <th>package</th>
    <th>說明</th>
    <th>操作說明</th>
    <tr>
        <td><a href="https://github.com/spf13/viper" target="_blank">viper</a></td>
        <td>Viper是一個配置設定文件、環境變量</td>
        <td>-</td>
    </tr>
     <tr>
        <td><a href="https://github.com/uber-go/zap" target="_blank">zap</a></td>
        <td>Zap 是一個快速、結構化、級別化的日誌庫，由 Uber 開發</td>
        <td> <a href="./asset/markdown/zap.md" target="_blank">open</a>  </td>
    </tr>
    <tr>
        <td><a href="https://github.com/gin-contrib/zap" target="_blank">gin zap middleware</a></td>
        <td>Gin 框架封裝的 zap 日誌中間件</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/lestrrat-go/file-rotatelogs" target="_blank">file-rotatelogs</a></td>
        <td>Go 語言的日誌文件切割和彙整庫</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/golang/crypto/tree/master" target="_blank">crypto/bcrypt</a></td>
        <td>字串加密核對</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/go-gorm/gorm" target="_blank">gorm</a></td>
        <td>Go 語言 ORM 庫，它支持 MySQL、PostgreSQL、SQLite 和 SQL Server 數據庫</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/go-sql-driver/mysql" target="_blank">go-sql-driver/mysql</a></td>
        <td>MySQL 驅動，連接 MySQL 數據庫</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/golang-jwt/jwt" target="_blank">golang-jwt</a></td>
        <td>JSON Web Token (JWT) 庫</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/go-playground/validator" target="_blank">validator</a></td>
        <td>驗證器用於驗證結構體和個別的數據</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/gin-contrib/cors" target="_blank">cors</a></td>
        <td>跨域請求的中間件</td>
        <td> - </td>
    </tr> 
    <tr>
        <td><a href="https://github.com/redis/go-redis/v9" target="_blank">go-redis</a></td>
        <td>go-redis 是 Redis 客户端库</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/swaggo/gin-swagger" target="_blank">gin-swagger</a></td>
        <td>gin swagger 產生 API docs</td>
        <td> <a href="./asset/markdown/swagger.md" target="_blank">open</a> </td>
    </tr>
    <tr>
        <td><a href="https://github.com/hibiken/asynq" target="_blank">asynq</a></td>
        <td>基於 Redis 的分布式任務隊列和異步處理庫，用於實現事件驅動架構</td>
        <td><a href="https://www.tizi365.com/topic/14001.html" target="_blank">open</a>  </td>
    </tr>
    <tr>
        <td><a href="https://github.com/prometheus/client_golang" target="_blank">prometheus client_golang</a></td>
        <td>Prometheus 客户端庫，收集 HTTP 請求、業務事件、系統指標</td>
        <td> - </td>
    </tr>
</table>