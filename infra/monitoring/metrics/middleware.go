package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 是一個 Gin 中間件，用於記錄 HTTP 請求指標
func MetricsMiddleware(collector *MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 處理請求
		c.Next()

		// 計算持續時間
		duration := time.Since(start).Seconds()

		// 獲取請求資訊
		method := c.Request.Method
		path := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())
		apiVersion := "v1" // 可從路徑中解析

		// 記錄指標
		collector.RecordHttpRequest(method, path, status, apiVersion, duration)
	}
}
