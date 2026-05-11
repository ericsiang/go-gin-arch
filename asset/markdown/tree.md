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
│   ├── admin                   => 後台管理員領域
│   │   ├── entity              => 存放所有的實體 (包含聚合根與內部實體)
│   │   │       └── admin.go
│   │   ├── repository          => 資料操作，負責使用 dao 進行資料操作
│   │   │   ├── dao             => 資料存取層
│   │   │   │   └── admin_dao.go
│   │   │   ├── model           => 資料表結構的 struct
│   │   │   │   └── admin.go
│   │   │   └── admin_repo.go
│   │   └── service             => 業務邏輯處理
│   │       └── admin_serv.go
│   └── user                    => 用戶領域
│       ├── entity              => 存放所有的實體 (包含聚合根與內部實體)
│       │       └── users.go    => 聚合根
│       ├── valueobj            => 值物件
│       ├── events              => 用戶事件處理
│       │   └── user_serv.go
│       ├── repository          => 資料操作，負責使用 dao 進行資料操作
│       │   ├── dao             => 資料存取層
│       │   │   └── user_dao.go
│       │   ├── model           => 資料表結構的 struct
│       │   │   └── user.go
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
│   │           ├── events      => 領域事件（Event-Driven）
│   │           │   └── user_events.go  => 用戶事件定義和處理器
│   │           └── users.go
│   ├── handler                 => 放置 gin 框架的 handler 程式碼
│   │   ├── handleError.go
│   │   ├── handleValidate.go
│   │   ├── handlerGeneric.go
│   │   └── handlerMysql.go
│   ├── inter                   => 放置 gin 框架內部使用的程式碼
│   │   └──response             => 放置 gin 框架內部使用的 response 程式碼
│   │        └── generic_resp.go
│   ├── middleware              => 放置 gin 框架的 middleware 程式碼
│   │   ├── jwt_auth.go
│   │   ├── opa_auth.go
│   │   └── rate_limit.go
│   ├── router                  => 放置 gin 框架的 router
│   │   └── router.go
│   └── validate_lang           => 放置 gin 框架的驗證語言設定
│       └── validate_lang.go
├── go.mod
├── go.sum
├── infra                       => 放置基礎建設的程式碼
│   ├── cache                   => 快取
│   │   └── redis
│   │       └── redis.go
│   ├── database                => 資料庫操作
│   │   ├── migrate             => 資料庫遷移
│   │   │   └── migrate.go
│   │   └── seeder              => 建立初始資料庫資料
│   │       ├── common_seeder.go
│   │       └── seeder.go
│   ├── env                     => 環境變數設定
│   │   ├── config.go
│   │   └── env.go
│   ├── event                   => 事件驅動架構基礎設施 
│   │   ├── event.go            => 事件核心定義和接口
│   │   ├── broker.go           => 事件代理器（EventBroker）
│   │   ├── asynq_client.go     => Asynq 客戶端實現（Publisher）
│   │   ├── asynq_server.go     => Asynq 服務端實現（Subscriber）
│   │   ├── rabbitmq_client.go  => RabbitMQ 客戶端實現（Publisher）
│   │   └── rabbitmq_server.go  => RabbitMQ 服務端實現（Subscriber）
│   ├── log                     => 日誌
│   │   └── zap_log
│   │       └── logger.go
│   └── orm                     => 資料庫 ORM
│       └── gorm_mysql
│           └── mysql.go
├── internal                    => 放置內部使用的程式碼，例如通用的 dao、model 等
│   ├── apperror
│   │   └── app_error.go        => 錯誤處理結構和相關函數
│   ├── dao
│   │   └── generic_gorm_dao.go => 提供通用資料存取功能，使用 GORM 作為 ORM 工具
│   └── model
│       └── gormModel.go        => 定義了 GORM 模型的基礎結構，供其他模型繼承使用
├── log                         => 置放 log 檔，可依需求將 log level 區分
│   ├── error
│   └── info
├── scripts                     => 各式腳本用資料夾
│   ├── docker                  => docker 建立容器的腳本
│   └── wrk                     => wrk 壓測用 lua 腳本
├── test                        => 放置測試用的程式碼
│   └── limit_ping_test.go
├── tree.md
├── tree_mvc.md
└── util                        => 置放封裝工具
    ├── bcryptEncap             => 字串加密核對
    │   ├── bcrypt.go
    │   └── bcryptEncap_test.go
    ├── jwt_secret              => jwt 操作
    │   ├── jwt_secret.go
    │   └── jwt_secret_test.go
    ├── mysql_manager
    │   └── mysql_err_code.go
    ├── open_policy_agent       => open policy agent 操作
    │   ├── rbac.go
    │   ├── rbac.rego
    │   └── rbac_test.rego
    ├── swagger_docs            => swagger docs 使用
    │   └── swag_params.go
    └── zap_logger              => zap plugin
        ├── zapLoggger_test.go
        └── zap_logger.go   

```