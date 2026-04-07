// Package database 定义数据库接口，提供连接、关闭和获取数据库实例的方法
package database

import "self_go_gin/infra/env"

// Database 定義資料庫接口，提供連接、關閉和獲取資料庫實例的方法
type Database interface {
	Connect(serverEnv *env.ServerConfig) error
	Close() error
	GetDB() interface{}
}
