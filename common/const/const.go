// Package constset 定義了整個應用程式中使用的常量
package constset

import "time"

const (
	// ShutdownTimeout 服务器關閉的超時時間
	ShutdownTimeout = 10 * time.Second
)
