// Package container 提供基础设施组件的初始化函数
package container

import (
	"fmt"

	"self_go_gin/infra/cache/redis"
	"self_go_gin/infra/env"
	"self_go_gin/infra/event"
	gormysql "self_go_gin/infra/orm/gorm_mysql"
	jwtsecret "self_go_gin/util/jwt_secret"

	validlang "self_go_gin/gin_application/validate_lang"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// InitMysql 初始化 MySQL 数据库连接
func InitMysql(config *env.ServerConfig) (*gorm.DB, error) {
	db, err := gormysql.InitMysql(config)

	if err != nil {
		return nil, fmt.Errorf("failed to get mysql db: %w", err)
	}
	fmt.Println("MySQL database initialized")
	return db, nil
}

// InitRedis 初始化 Redis 连接
func InitRedis(config *env.ServerConfig) (*redis.Client, error) {
	redisClient, err := goredis.InitRedis(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis: %w", err)
	}
	// 根据您的 redis 包实现，可能需要调整返回类型
	fmt.Println("Redis initialized")
	return redisClient, nil
}

// InitEventBroker 初始化事件代理
func InitEventBroker(config *env.ServerConfig) (*event.Broker, error) {
	broker, err := event.NewBroker(event.BrokerTypeAsynq, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create event broker: %w", err)
	}
	return broker, nil
}

// InitJWTSecret 初始化 JWT 密钥
func InitJWTSecret(secret string) {
	jwtsecret.SetJwtSecret(secret)
	fmt.Println("JWT secret configured")
}

// InitValidator 初始化验证器本地化
func InitValidator(lang string) error {
	if err := validlang.InitValidateLang(lang); err != nil {
		return fmt.Errorf("failed to init validator: %w", err)
	}
	fmt.Println("Validator localization initialized")
	return nil
}

// InitCommonComponents 初始化常用组件（JWT、验证器等）
func InitCommonComponents(config *env.ServerConfig) error {
	// 设置 JWT 密钥
	InitJWTSecret(config.JwtSecret)

	// 初始化验证器中文化
	if err := InitValidator("zh"); err != nil {
		return err
	}

	return nil
}
