# self-go-gin (golang gin framwork 設計自用模板)

### 檔案結構 (tree 指令產生)
```
.
├── README.md                   => 說明檔
├── asset                       => 放置素材檔案
├── cmd                         => 放置執行檔案
├── common                      => 放置通用宣告
│   ├── const                   => 設定常數
│   │   └── const.go
│   └── msg_id                  => 統一定義訊息識別碼
│       └── msg_id.go
├── conf                        => 放置環境變數設定檔案
│   ├── env.docker.yaml.example
│   └── env.yaml.example
├── container                   => 依賴注入容器層，統一管理所有組件的生命周期和依賴關係
│   ├── app_container.go        => 應用容器核心（單例模式、線程安全）
│   ├── infra_provider.go       => 基礎設施提供者（DB、Redis、EventBroker 等）
│   └── README.md               => 容器層使用文檔
├── domains                     => 放置 domain 層的程式碼，依據功能分為不同的子目錄
│   ├── admin                   => 後台管理員
│   │   ├── entity              => 資料模型
│   │   │   └── model           => 資料表結構的 struct
│   │   │       └── admin.go
│   │   ├── repository          => 資料操作，負責使用 dao 進行資料操作
│   │   │   ├── dao             => 資料存取層
│   │   │   │   └── admin_dao.go
│   │   │   └── admin_repo.go
│   │   └── service             => 業務邏輯處理
│   │       └── admin_serv.go
│   └── user                    => 用戶
│       ├── entity              => 資料模型
│       │   └── model           => 資料表結構的 struct
│       │       └── users.go
│       ├── events              => 用戶事件處理
│       │   └── user_serv.go
│       ├── repository          => 資料操作，負責使用 dao 進行資料操作
│       │   ├── dao             => 資料存取層
│       │   │   └── user_dao.go
│       │   └── user_repo.go
│       └── service             => 業務邏輯處理
│           └── user_serv.go
├── gin_application             => 放置 gin 框架的程式碼
│   ├── api                     => 放置 gin 框架的 api controller 程式碼
│   │   └── v1
│   │       ├── admin
│   │       │   ├── request
│   │       │   │   └── admin_req.go
│   │       │   ├── response
│   │       │   │   └── admin_resp.go
│   │       │   └── admin.go
│   │       └── user
│   │           ├── request
│   │           │   └── user_req.go
│   │           ├── response
│   │           │   └── user_resp.go
│   │           ├── events     => 領域事件（Event-Driven）
│   │           │   └── user_events.go  => 用戶事件定義和處理器
│   │           └── users.go
│   ├── handler                => 放置 gin 框架的 handler 程式碼
│   │   ├── handleError.go
│   │   ├── handleValidate.go
│   │   ├── handlerGeneric.go
│   │   └── handlerMysql.go
│   ├── inter                  => 放置 gin 框架內部使用的程式碼
│   │   └──response            => 放置 gin 框架內部使用的 response 程式碼
│   │        └── generic_resp.go
│   ├── middleware             => 放置 gin 框架的 middleware 程式碼
│   │   ├── jwt_auth.go
│   │   ├── opa_auth.go
│   │   └── rate_limit.go
│   ├── router                 => 放置 gin 框架的 router
│   │   └── router.go
│   └── validate_lang          => 放置 gin 框架的驗證語言設定
│       └── validate_lang.go
├── go.mod
├── go.sum
├── infra                      => 放置基礎建設的程式碼
│   ├── cache                  => 快取
│   │   └── redis
│   │       └── redis.go
│   ├── database               => 資料庫操作
│   │   ├── migrate            => 資料庫遷移
│   │   │   └── migrate.go
│   │   └── seeder             => 建立初始資料庫資料
│   │       ├── common_seeder.go
│   │       └── seeder.go
│   ├── env                    => 環境變數設定
│   │   ├── config.go
│   │   └── env.go
│   ├── event                   => 事件驅動架構基礎設施 
│   │   ├── event.go            => 事件核心定義和接口
│   │   ├── broker.go           => 事件代理器（EventBroker）
│   │   ├── asynq_client.go     => Asynq 客戶端實現（Publisher）
│   │   └── asynq_server.go     => Asynq 服務端實現（Subscriber）
│   ├── log                     => 日誌
│   │   └── zap_log
│   │       └── logger.go
│   └── orm                     => 資料庫 ORM
│       └── gorm_mysql
│           └── mysql.go
├── internal                    => 放置內部使用的程式碼，例如通用的 dao、model 等
│   ├── dao
│   │   └── generic_dao.go
│   └── model
│       ├── gormModel.go
│       └── model_setting.go
├── log                        => 置放 log 檔，可依需求將 log level 區分
│   ├── error
│   └── info
├── scripts                    => 各式腳本用資料夾
│   └── docker                 => docker 建立容器的腳本
├── test                       => 放置測試用的程式碼
│   └── limit_ping_test.go
├── tree.md
├── tree_mvc.md
└── util                       => 置放封裝工具
    ├── bcryptEncap            => 字串加密核對
    │   ├── bcrypt.go
    │   └── bcryptEncap_test.go
    ├── jwt_secret             => jwt 操作
    │   ├── jwt_secret.go
    │   └── jwt_secret_test.go
    ├── mysql_manager
    │   └── mysql_err_code.go
    ├── open_policy_agent      => open policy agent 操作
    │   ├── rbac.go
    │   ├── rbac.rego
    │   └── rbac_test.rego
    ├── swagger_docs            => swagger docs 使用
    │   └── swag_params.go
    ├── track_time              => 計算 func 程式時間
    │   ├── track_time.go
    │   └── track_time_test.go
    └── zap_logger              => zap plugin
        ├── zapLoggger_test.go
        └── zap_logger.go   

```
### 專案介紹
#### 這是一個基於 Go 語言開發的後端 web service 模板，旨在提供一個結構清晰、易於擴展和維護的代碼基礎，目前是搭配 Gin 框架構建，此結構有助於未來替換 Web 框架（例如從 Gin 換成 Echo），降低替換成本
* 分層架構
  * 初始化容器層 (Container)
      * 統一管理所有依賴實例（DB、Redis、EventBroker 等）
      * 單例模式保證全局唯一容器實例
      * 線程安全（使用 sync.RWMutex）
    * 生命周期管理
      * `Initialize()` - 統一初始化所有基礎設施
      * `Shutdown()` - 優雅關閉所有資源
    * 基礎設施組件初始化（MySQL、Redis、EventBroker、JWT、Validator）
    * 優勢
      * 解耦：各層通過容器獲取依賴，避免直接依賴
      * 易維護：依賴管理集中，修改配置只需調整容器層
      * 可擴展：輕鬆添加新的基礎設施組件

  * 採用 DDD (Domain-Driven Design) 思維，清晰的職責分離與職權域邊界
  * 完整的分層架構設計 (Entity → ValueObject → Service → Repository → DAO → Service)
    * Entity 層（聚合根）**
      * 建立實體聚合根，包含身份與狀態
      * ValueObject 驗證：內置值物件進行數據驗證和加密（如 Account、Password）
    * Service 層
      * 應用層邏輯協調（事務邊界、流程編排、複雜業務流程）
      * 調用 Repository 進行數據操作
    * Repository 層
      * 領域模型（Domain Model）與持久化模型（PO）之間的轉換
      * cache 的調用也在這層
      * 隱藏所有數據庫實現細節，對上層提供純粹的領域模型
      * DAO 層
        * 純粹的數據庫操作，直接使用 GORM 和 PO（持久化物件）
        * 不了解業務邏輯，只執行 SQL 操作
  * 符合關注點分離原則
  * 可維護性高，修改業務邏輯只需動 domains 資料夾
  
* 基礎設施
  * 環境配置管理(Viper)
  * 日誌系統 (Zap)
  * 快取機制 (Redis)
  * 資料庫連線 (Mysql)
  * 資料庫遷移和種子資料 (migrate、seeder)
  * 事件驅動架構 (Event-Driven 透過環境設定檔控制是否使用) [架構說明](./asset/markdown/event-driven.md)

* 安全性考量
  * JWT 認證
  * OPA 權限控制
  * Bcrypt 加密核對
  
* Web 框架 (gin_application)
  * router
  * 中間件
    * 全域追蹤機制 (Trace-ID)
    * 限流機制
    * JWT 認證機制
    * 權限驗證機制
  * 統一的 response 處理
  * API 版本控制
  
* 標準化與規範的開發實踐
  * 統一的錯誤處理
  * 參數驗證機制
  * Swagger 文檔支援
  * 測試檔案配置 
  * gin 框架相關程式碼集中於 /gin_application 
  * 可擴展性高，可輕鬆添加新的功能模組

* 可觀測性設計  
  * 全域追蹤 (Tracing)： Middleware 自動為每個 Request 生成 Trace-ID 並注入 Context，確保跨層級 Log 關聯，方便未來追蹤
  * 結構化日誌 (Structured Logging)： 整合 Zap Logger 並區分 Info 與 Error 級別
  
* 優化功能
  * Graceful Shutdown： 停止接收 request，timeout 內等待已接收連線處理結束
  
* 容器化部署
  * 透過 docker 快速建立容器

### 架構圖結構
#### HTTP 請求處理流程 (DDD 分層架構)
``` mermaid
graph LR
    Client((用戶端)) -->|HTTP Request| Router["Router<br/>路由匹配"]
    
    Router --> Middleware["Middleware<br/>日誌/認證/授權"]
    
    Middleware --> Controller["Controller<br/>參數驗證與轉換"]
    
    Controller --> Service["ApplicationService<br/>事務/流程協調/業務規則檢查"]
    
    Service --> Repository["Repository<br/>領域模型 ↔ PO 轉換"]
    
    Repository --> DAO["DAO<br/>純粹數據庫操作"]
    
    DAO --> Database[("Database<br/>MySQL/Redis")]
    
    Database --> DAO
    DAO --> Repository
    Repository --> Service
    Service --> Controller
    Controller --> Response["Response<br/>統一格式輸出"]
    Response --> Client
    
```

#### 框架可替換性設計
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

#### 事件驅動架構流程
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



#### 初始化容器層架構
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

### 快速开始
* conf 資料夾內，依需求複製配置文件範例，檔名不含.example，並填入真實配置
* 查看執行 make file

### 使用到的 package
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
        <td> - </td>
    </tr>
</table>