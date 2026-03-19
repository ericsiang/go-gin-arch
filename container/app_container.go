// Package container 提供依赖注入容器，管理应用中所有组件的生命周期和依赖关系
package container

import (
	"context"
	"fmt"
	"sync"

	"self_go_gin/infra/env"
	"self_go_gin/infra/event"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AppContainer 应用容器，负责管理所有依赖实例
type AppContainer struct {
	mu sync.RWMutex
	// 環境配置
	config *env.ServerConfig
	// 基础设施
	db          *gorm.DB
	redisClient *redis.Client
	eventBroker *event.Broker
	// 其他通用实例可以在此添加
	// 例如: logger, metrics, tracer 等
}

var (
	instance *AppContainer
	once     sync.Once
)

// GetContainer 获取全局容器实例（单例模式）
func GetContainer() *AppContainer {
	once.Do(func() {
		instance = &AppContainer{}
	})
	return instance
}

// Initialize 初始化容器及所有依赖
func (c *AppContainer) Initialize(configPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("\n Initializing application container...")

	// 1. 加载配置
	if err := env.InitEnv(configPath); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	c.config = env.GetConfigManager().GetServerEnv()
	fmt.Println("Configuration loaded")

	// 2. 设置 Gin 模式
	gin.SetMode(c.config.AppMode)
	fmt.Printf("Gin mode: %s\n", c.config.AppMode)

	// 3. 初始化基础设施
	if err := c.initInfrastructure(); err != nil {
		return fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	fmt.Println("\n All components initialized successfully!")
	return nil
}

// initInfrastructure 初始化所有基础设施组件
func (c *AppContainer) initInfrastructure() error {
	// 初始化数据库
	if c.config.MysqlDB.IsEnabled {
		db, err := InitMysql(c.config)
		if err != nil {
			return fmt.Errorf("failed to init mysql: %w", err)
		}
		if db == nil {
			return fmt.Errorf("mysql db instance is nil")
		}
		c.db = db
	} else {
		fmt.Println("MySQL is not enabled by env")
	}

	// 初始化 Redis
	if c.config.Redis.IsEnabled {
		redisClient, err := InitRedis(c.config)
		if err != nil {
			return fmt.Errorf("failed to init redis: %w", err)
		}
		if redisClient == nil {
			return fmt.Errorf("redis client instance is nil")
		}
		c.redisClient = redisClient
	} else {
		fmt.Println("Redis is not enabled by env")
	}

	// 初始化事件代理
	if c.config.IsEventBroker {
		broker, err := InitEventBroker(c.config)
		if err != nil {
			return fmt.Errorf("failed to init event broker: %w", err)
		}
		if broker == nil {
			return fmt.Errorf("event broker instance is nil")
		}
		c.eventBroker = broker
		fmt.Println("Event broker ready (Publisher only)")
		fmt.Println("Note: To process events, run the event_worker service")
	} else {
		fmt.Println("EventBroker is not enabled by env")
	}

	return nil
}

// Shutdown 优雅关闭所有资源
func (c *AppContainer) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("Shutting down application container...")

	var errs []error

	// 关闭事件代理
	if c.eventBroker != nil {
		if err := c.eventBroker.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close event broker: %w", err))
		}
	}

	// 关闭 Redis 连接
	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close redis client: %w", err))
		} else {
			fmt.Println("Redis client closed")
		}
	}

	// 关闭数据库连接
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close database: %w", err))
			} else {
				fmt.Println("Database connection closed")
			}
		}
	}

	// 可以在此添加其他资源清理逻辑

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	fmt.Println("Application container shutdown completed")
	return nil
}

// GetDB 获取数据库实例
func (c *AppContainer) GetDB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

// GetRedisClient 获取 Redis 客户端实例
func (c *AppContainer) GetRedisClient() *redis.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redisClient
}

// GetEventBroker 获取事件代理实例
func (c *AppContainer) GetEventBroker() *event.Broker {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.eventBroker
}

// GetConfig 获取配置实例
func (c *AppContainer) GetConfig() *env.ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
