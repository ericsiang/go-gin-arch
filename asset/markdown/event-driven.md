## 事件驅動架構 

### 🎯 設計目標
將事件驅動架構設計為接口模式，實現底層實現的可插拔性，方便未來從 Asynq 切換到 RabbitMQ、Kafka 或其他消息隊列。

* 基於發佈-訂閱模式 (Pub/Sub Pattern)
* 支持多種事件代理（目前實現 Asynq，未來可擴展替換為 RabbitMQ、Kafka）
* 領域事件與領域邏輯解耦，提高系統可擴展性

### 📦 核心組件

#### EventBroker - 事件代理器
統一管理事件發佈與訂閱，封裝 Publisher 和 Subscriber。

**主要方法：**
- `Publisher()` - 獲取事件發佈者實例
- `Subscriber()` - 獲取事件訂閱者實例
- `BrokerType()` - 獲取代理類型（asynq/rabbitmq/kafka）
- `Close()` - 關閉事件代理

#### Publisher - 事件發佈者接口
負責發佈事件到消息隊列，支持同步/異步/延遲發佈。

**核心方法：**
- `Publish(ctx, event)` - 使用默認選項發佈事件
- `PublishWithOptions(ctx, event, opts)` - 使用自定義選項發佈事件
- `PublishDeferred(ctx, event, delaySeconds)` - 延遲發佈事件
- `Close()` - 關閉發佈者

#### Subscriber - 事件訂閱者接口
負責事件消費和處理，管理事件處理器。

**核心方法：**
- `Subscribe(handler)` - 訂閱事件，註冊處理器
- `Start()` - 啟動訂閱者（非阻塞）
- `Run()` - 運行訂閱者（阻塞）
- `Shutdown()` - 優雅關閉訂閱者

#### Handler - 事件處理器接口
各領域實現自己的事件處理邏輯。

**必須實現的方法：**
- `EventType()` - 返回該處理器處理的事件類型
- `Handle(ctx, event)` - 處理事件的核心邏輯

### ⚙️ 功能特性
* **事件優先級設定** - 支持高、中、低優先級隊列
* **自動重試機制** - 可配置重試次數（默認 3 次）
* **延遲發佈支持** - 實現定時任務功能
* **分布式追蹤** - 支持 TraceID 用於分布式追蹤
* **優雅關閉** - 確保所有事件處理完成後才關閉
* **並發處理** - 支持多協程並發處理事件
* **隊列管理** - 支持多隊列配置和優先級權重

### 🎬 應用場景
* **用戶註冊流程** - 註冊成功後發送歡迎郵件、記錄審計日誌
* **資料變更同步** - 資料更新後同步至其他服務或緩存
* **異步處理耗時任務** - 圖片處理、報表生成、文件導出
* **業務解耦** - 分離核心業務與輔助功能，提升系統效能
* **定時任務** - 通過延遲發佈實現定時提醒、定期任務

### ✨ 優勢
* **領域解耦** - 領域間通過事件通信，不需直接依賴
* **易於擴展** - 新增事件處理器不影響現有功能
* **效能優化** - 異步處理不阻塞主流程，提升用戶體驗
* **可測試性** - 事件處理器可獨立測試
* **可替換性** - 基於接口設計，可輕鬆切換底層實現

### 架構設計

#### 📂 文件結構
```
infra/event/
├── event.go          # 核心接口定義（Event, Handler, Publisher, Subscriber）
├── broker.go         # 事件代理工廠和封裝
├── asynq_client.go   # Asynq 發布者實現（基於 Redis）
└── asynq_server.go   # Asynq 訂閱者實現（基於 Redis）
```

**領域事件目錄：**
```
domains/user/
└── events/
    └── user_events.go  # 用戶領域的事件定義和處理器
```

**事件處理服務：**
```
cmd/
├── event_worker/
│   └── main.go  # 獨立的事件處理服務入口
└── first_web_service/
    └── main.go  # Web 服務（發布事件）
```

#### 🔧 核心數據結構

**Event 事件結構：**
```go
type Event struct {
    ID        string              // 事件唯一標識
    Type      string              // 事件類型（如 "user.created"）
    Payload   json.RawMessage     // 事件負載數據
    Source    string              // 事件來源
    TraceID   string              // 分布式追踪 ID
    Timestamp time.Time           // 事件時間戳
    Metadata  map[string]string   // 元數據
}
```

**PublishOptions 發布選項：**
```go
type PublishOptions struct {
    Queue    string          // 隊列名稱（default/high/low）
    Priority int             // 優先級（1-10，10 最高）
    MaxRetry int             // 最大重試次數
    Delay    time.Duration   // 延遲時間
    Timeout  time.Duration   // 處理超時時間
}
```

#### 🎯 Asynq 實現細節

**隊列配置：**
- **HighPriorityQueue（high）** - 權重 6，處理緊急任務
- **DefaultQueue（default）** - 權重 3，處理普通任務
- **LowPriorityQueue（low）** - 權重 1，處理低優先級任務

**並發配置：**
- **Concurrency: 10** - 同時處理 10 個任務
- **自動重試** - 默認重試 3 次，可自定義
- **錯誤處理** - 自動記錄處理失敗的任務

**Redis 配置：**
```yaml
# config/env.yaml
IsEventBroker: true
Redis:
  IsEnabled: true
  Host: "localhost"
  Port: 6379
  Password: ""
  DBName: 0
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



### 事件處理器註冊流程
``` mermaid
sequenceDiagram
    participant Main as cmd/event_worker<br/>main.go
    participant Container as Container<br/>容器層
    participant Broker as EventBroker<br/>事件代理
    participant Sub as Subscriber<br/>訂閱者
    participant Handler as EventHandler<br/>事件處理器
    
    Main->>Container: 1. Initialize(configPath)
    Container->>Container: 2. 初始化 Redis
    Container->>Broker: 3. NewEventBroker()
    Broker-->>Container: 返回 EventBroker
    Container-->>Main: 完成初始化
    
    Main->>Broker: 4. GetSubscriber()
    Broker-->>Main: 返回 Subscriber
    
    Main->>Handler: 5. NewUserCreatedHandler()
    Handler-->>Main: 返回 Handler
    
    Main->>Sub: 6. Subscribe(Handler)
    Sub->>Sub: 註冊事件類型與處理器映射
    Sub-->>Main: 註冊成功
    
    Main->>Sub: 7. Run()
    Sub->>Sub: 啟動工作協程池
    Sub->>Sub: 開始監聽 Redis 隊列
    
    Note over Sub: 等待事件到來...
    
```

---

## 📚 使用指南

### 1️⃣ 定義領域事件

**步驟 1：定義事件類型常量**
```go
// domains/user/events/user_events.go
const (
    UserCreatedEventType = "user.created"
    UserUpdatedEventType = "user.updated"
    UserDeletedEventType = "user.deleted"
)
```

**步驟 2：定義事件負載結構**
```go
// UserCreatedEventPayload 用戶創建事件的負載
type UserCreatedEventPayload struct {
    UserID   uint   `json:"user_id"`
    Account  string `json:"account"`
    Email    string `json:"email,omitempty"`
    CreateAt string `json:"create_at"`
}
```

**步驟 3：實現事件處理器**
```go
// UserCreatedEventHandler 處理用戶創建事件
type UserCreatedEventHandler struct{}

func NewUserCreatedEventHandler() *UserCreatedEventHandler {
    return &UserCreatedEventHandler{}
}

// EventType 返回處理的事件類型
func (h *UserCreatedEventHandler) EventType() string {
    return UserCreatedEventType
}

// Handle 處理用戶創建事件
func (h *UserCreatedEventHandler) Handle(ctx context.Context, evt *event.Event) error {
    var payload UserCreatedEventPayload
    if err := evt.UnmarshalPayload(&payload); err != nil {
        return fmt.Errorf("failed to unmarshal payload: %w", err)
    }

    // 實現業務邏輯
    log.Printf("Processing user creation: UserID=%d, Account=%s", 
        payload.UserID, payload.Account)
    
    // 發送歡迎郵件
    if err := h.sendWelcomeEmail(payload); err != nil {
        log.Printf("Failed to send welcome email: %v", err)
        // 可選：返回錯誤以觸發重試
    }
    
    // 記錄審計日誌
    h.logAudit(payload)
    
    return nil
}
```

### 2️⃣ 在業務服務中發布事件

**在 Service 層發布事件：**
```go
// domains/user/service/user_serv.go
type UserService struct {
    repo      repository.UserRepository
    publisher event.Publisher  // 注入 Publisher
}

// CreateUser 創建用戶並發布事件
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // 1. 執行核心業務邏輯
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    // 2. 發布領域事件（異步處理）
    if err := s.publishUserCreatedEvent(ctx, user); err != nil {
        // 記錄錯誤但不影響主流程
        log.Printf("Failed to publish user created event: %v", err)
    }
    
    return nil
}

// publishUserCreatedEvent 發布用戶創建事件
func (s *UserService) publishUserCreatedEvent(ctx context.Context, user *model.User) error {
    payload := events.UserCreatedEventPayload{
        UserID:   user.ID,
        Account:  user.GetAccount(),
        CreateAt: time.Now().Format(time.RFC3339),
    }
    
    evt, err := event.NewEvent(events.UserCreatedEventType, payload)
    if err != nil {
        return fmt.Errorf("failed to create event: %w", err)
    }
    
    evt.Source = "user-service"
    
    if err := s.publisher.Publish(ctx, evt); err != nil {
        return fmt.Errorf("failed to publish event: %w", err)
    }
    
    log.Printf("[UserService] User created event published: UserID=%d", user.ID)
    return nil
}
```

**獲取 Publisher 實例：**
```go
// 方式 1：從容器獲取
app := container.GetContainer()
broker := app.GetEventBroker()
publisher := broker.Publisher()

// 方式 2：直接獲取 Asynq Client
publisher := event.GetAsynqClient()
```

### 3️⃣ 註冊事件處理器

**在 cmd/event_worker/main.go 中註冊：**
```go
func main() {
    // 初始化容器
    app := container.GetContainer()
    app.Initialize(configPath)
    
    // 獲取事件代理
    broker := app.GetEventBroker()
    subscriber := broker.Subscriber()
    
    // 註冊事件處理器
    if err := registerEventHandlers(subscriber); err != nil {
        log.Fatalf("Failed to register handlers: %v", err)
    }
    
    // 啟動事件處理服務
    subscriber.Run()
}

func registerEventHandlers(subscriber event.Subscriber) error {
    handlers := []event.Handler{
        events.NewUserCreatedEventHandler(),
        events.NewUserUpdatedEventHandler(),
        events.NewUserDeletedEventHandler(),
    }
    
    for _, handler := range handlers {
        if err := subscriber.Subscribe(handler); err != nil {
            return fmt.Errorf("failed to subscribe handler %s: %w", 
                handler.EventType(), err)
        }
    }
    
    log.Printf("Successfully registered %d handlers", len(handlers))
    return nil
}
```

### 4️⃣ 高級用法

#### 延遲發布事件
```go
// 30 秒後發布事件
err := publisher.PublishDeferred(ctx, evt, 30)
```

#### 使用自定義選項發布
```go
opts := &event.PublishOptions{
    Queue:    event.HighPriorityQueue,  // 使用高優先級隊列
    Priority: 10,                        // 最高優先級
    MaxRetry: 5,                         // 重試 5 次
    Delay:    time.Minute * 5,          // 延遲 5 分鐘
    Timeout:  time.Second * 30,         // 超時 30 秒
}
err := publisher.PublishWithOptions(ctx, evt, opts)
```

#### 發布到特定隊列
```go
// 高優先級隊列
client := event.GetAsynqClient()
err := client.PublishToHighPriorityQueue(ctx, evt)

// 低優先級隊列
err := client.PublishToLowPriorityQueue(ctx, evt)
```

---

## 🔍 最佳實踐

### ✅ 事件設計原則
1. **事件名稱規範** - 使用 `{domain}.{action}` 格式（如 `user.created`）
2. **事件不可變** - 事件發布後不應修改
3. **包含完整信息** - Payload 應包含處理所需的所有數據
4. **冪等性設計** - 事件處理器應支持重複處理
5. **明確事件來源** - 設置 `Source` 字段標識事件來源

### ✅ 錯誤處理策略
1. **區分可重試錯誤** - 網絡錯誤、暫時性失敗應重試
2. **記錄處理失敗** - 使用日誌記錄所有處理失敗的事件
3. **設置合理超時** - 避免長時間阻塞
4. **死信隊列** - Asynq 自動處理失敗任務，可查看 Redis 中的死信

### ✅ 性能優化
1. **批量處理** - 處理器內部可以實現批量操作
2. **並發控制** - 合理配置 `Concurrency` 參數
3. **隊列分離** - 不同類型的任務使用不同隊列
4. **監控告警** - 監控隊列長度、處理延遲、失敗率

### ✅ 測試建議
1. **單元測試** - 測試事件處理器的業務邏輯
2. **集成測試** - 測試事件發布和消費的完整流程
3. **Mock Publisher** - 測試時可以 Mock Publisher避免實際發布

```go
// 測試示例
func TestUserCreatedEventHandler(t *testing.T) {
    handler := events.NewUserCreatedEventHandler()
    
    payload := events.UserCreatedEventPayload{
        UserID:  1,
        Account: "testuser",
    }
    
    evt, _ := event.NewEvent(events.UserCreatedEventType, payload)
    
    err := handler.Handle(context.Background(), evt)
    assert.NoError(t, err)
}
```

---

## 🚀 運行命令

### 啟動 Web 服務（發布事件）
```bash
# 設置配置路徑
export CONFIG_PATH=./conf/

# 啟動 Web 服務
go run cmd/first_web_service/main.go
```

### 啟動事件處理服務（消費事件）
```bash
# 設置配置路徑
export CONFIG_PATH=./conf/

# 啟動事件處理服務
go run cmd/event_worker/main.go
```

### 監控 Redis 隊列
```bash
# 連接 Redis
redis-cli

# 查看所有隊列
KEYS asynq:*

# 查看隊列長度
LLEN asynq:queues:default
LLEN asynq:queues:high
LLEN asynq:queues:low

# 查看進行中的任務
ZCARD asynq:active:default

# 查看待處理任務
ZCARD asynq:pending:default

# 查看失敗任務
ZCARD asynq:dead
```

---

## 🔄 未來擴展

### 切換到其他消息隊列

由於採用接口設計，切換到其他消息隊列非常簡單：

**1. 實現新的 Publisher 和 Subscriber**
```go
// infra/event/rabbitmq_client.go
type RabbitMQClient struct {
    // ...
}

func (c *RabbitMQClient) Publish(ctx context.Context, event *Event) error {
    // 實現 RabbitMQ 發布邏輯
}

// infra/event/rabbitmq_server.go
type RabbitMQServer struct {
    // ...
}

func (s *RabbitMQServer) Subscribe(handler Handler) error {
    // 實現 RabbitMQ 訂閱邏輯
}
```

**2. 在 broker.go 中添加新的實現**
```go
case BrokerTypeRabbitMQ:
    publisher = InitRabbitMQClient(config)
    subscriber = InitRabbitMQServer(config)
```

**3. 更新配置即可切換**
無需修改業務代碼，只需在容器層配置使用的 BrokerType。

---
## 📖 參考資源
- [Asynq Documentation](https://github.com/hibiken/asynq)
---

## 💡 常見問題

**Q: 事件處理失敗會怎樣？**
A: Asynq 會自動重試，默認重試 3 次。重試失敗後任務會進入死信隊列，可手動處理。

**Q: 如何保證事件處理的順序？**
A: 同一隊列的事件按 FIFO 順序處理。如需嚴格順序，使用單一隊列和並發數為 1。

**Q: 事件處理器可以發布新事件嗎？**
A: 可以。事件處理器內部可以發布新事件，形成事件鏈。

**Q: 如何避免事件重複處理？**
A: 實現冪等性設計，使用唯一標識符（如事件 ID）去重，或在業務層面保證冪等。

**Q: Web 服務和事件處理服務必須分離嗎？**  
A: 推薦分離部署以提高可擴展性和容錯能力，但也可以在同一進程中運行。