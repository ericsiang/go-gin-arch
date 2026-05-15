// Package constset 定義了整個應用程式中使用的常量
package constset

import "time"

const (
	// ShutdownTimeout 服务器關閉的超時時間
	ShutdownTimeout = 10 * time.Second
	// RedisRateLimitSecond 	限流的時間窗口（秒）
	RedisRateLimitSecond = 5
	// MySQLMaxIdleConns MySQL 連接池的最大空閒連接數
	MySQLMaxIdleConns = 100
	// MySQLMaxOpenConns MySQL 連接池的最大打開連接數
	MySQLMaxOpenConns = 150
)
