package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceMiddleware 請求 ID 中間件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 嘗試從 Header 獲取客戶端傳入的 Request ID
		traceID := c.GetHeader("X-Request-ID")

		// 如果客戶端沒有傳入，則自動生成
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 設置到 Context 中，供後續處理使用
		c.Set("trace_id", traceID)

		// 返回給客戶端（方便追蹤）
		c.Header("X-Request-ID", traceID)

		c.Next()
	}
}
