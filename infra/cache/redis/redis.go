// Package goredis 提供 Redis 客戶端的初始化和訪問功能，使用 go-redis 庫來與 Redis 服務器進行通信。
package goredis

import (
	"context"
	"fmt"
	"os"
	"self_go_gin/infra/env"
	"strconv"

	"github.com/redis/go-redis/v9"
)



// InitRedis 初始化 Redis 客戶端
func InitRedis(serverEnv *env.ServerConfig) (*redis.Client,error) {
	redisConfig := serverEnv.Redis
	redisAddr := redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisConfig.Password,
		DB:       0, // use default DB
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Fprintln(os.Stderr, "redis connect failed, err:", err)
		return nil, err
	}

	fmt.Println("redis client connect ping success")

	return redisClient, nil
}

