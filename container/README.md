# Container 容器层

## 概述
Container 层提供依赖注入容器，统一管理应用中所有组件的生命周期和依赖关系。

## 结构

```
container/
├── app_container.go      # 应用容器主体（单例模式，线程安全）
├── infra_provider.go     # 基础设施提供者（DB, Redis, EventBroker 等）
└── handler_provider.go   # HTTP 处理器提供者
```

## 使用示例

### 在 main.go 中初始化容器

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "self_go_gin/container"
    "self_go_gin/gin_application/router"
    "syscall"
    "time"
)

func main() {
    // 1. 获取配置路径
    configPath := os.Getenv("CONFIG_PATH")
    
    // 2. 初始化容器
    app := container.GetContainer()
    if err := app.Initialize(configPath); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize container: %v\n", err)
        os.Exit(1)
    }
    
    // 3. 初始化通用组件（JWT、验证器等）
    if err := container.InitCommonComponents(app.GetConfig()); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize common components: %v\n", err)
        os.Exit(1)
    }
    
    // 4. 启动 HTTP 服务
    httpServerRun(app)
}

func httpServerRun(app *container.AppContainer) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    // 设置路由
    router := router.Router(quit)
    config := app.GetConfig()
    
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", config.Port),
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // 启动服务器
    go func() {
        fmt.Printf("Server listening on :%d\n", config.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
            os.Exit(1)
        }
    }()
    
    // 等待退出信号
    <-quit
    fmt.Println("\nShutting down gracefully...")
    
    // 关闭 HTTP 服务器
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    
    // 关闭容器资源
    app.Shutdown()
    
    fmt.Println("Server exited")
}
```

### 在其他地方获取依赖实例

```go
package somepackage

import "self_go_gin/container"

func SomeFunction() {
    // 获取容器实例
    app := container.GetContainer()
    
    // 获取数据库连接
    db := app.GetDB()
    
    // 获取事件代理
    broker := app.GetEventBroker()
    
    // 获取配置
    config := app.GetConfig()
    
    // 使用这些实例...
}
```

### 使用 HandlerProvider 管理处理器

```go
package main

import "self_go_gin/container"

func setupHandlers() {
    app := container.GetContainer()
    handlerProvider := container.NewHandlerProvider(app)
    
    // 可以添加方法来获取各种 Handler
    // userHandler := handlerProvider.GetUserHandler()
    // adminHandler := handlerProvider.GetAdminHandler()
}
```

## 特性

- **单例模式**: 全局唯一容器实例
- **线程安全**: 使用 sync.RWMutex 保护共享资源
- **生命周期管理**: 统一初始化和关闭所有依赖
- **依赖隔离**: 各层通过容器获取依赖，降低耦合
- **可扩展**: 轻松添加新的基础设施或处理器提供者

## 扩展方式

### 添加新的基础设施组件

在 `infra_provider.go` 中添加：

```go
func InitLogger(config *env.ServerEnv) (*Logger, error) {
    // 初始化逻辑
}
```

在 `app_container.go` 中：

```go
type AppContainer struct {
    // ... 现有字段
    logger *Logger
}

func (c *AppContainer) GetLogger() *Logger {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.logger
}
```

### 添加新的 Handler 提供者

在 `handler_provider.go` 中添加：

```go
func (p *HandlerProvider) GetUserHandler() *UserHandler {
    db := p.container.GetDB()
    return NewUserHandler(db)
}
```
