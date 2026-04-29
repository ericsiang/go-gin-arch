package middleware

import (
	"self_go_gin/container"
	"self_go_gin/infra/monitoring/metrics"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 是一個 Gin 中間件，用於記錄 HTTP 請求指標
func MetricsMiddleware() gin.HandlerFunc {
	appName := container.GetContainer().GetConfig().AppName
	collector := metrics.NewMetricsCollector(appName)
	return func(c *gin.Context) {
		start := time.Now()

		// 處理請求
		c.Next()

		// 獲取請求資訊
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			// 處理 404 或未定義路由
			if c.Writer.Status() == 404 {
				path = "not_found"
			} else {
				path = "root"
			}
		}
		status := strconv.Itoa(c.Writer.Status())
		apiVersion := strings.SplitN(path, "/", 2)[1] // api/v1/resource

		// 計算持續時間
		duration := time.Since(start).Seconds()

		// 記錄指標
		collector.RecordHTTPRequest(method, path, status, apiVersion, duration)
	}
}
