// Package ginlogger 提供 Zap 日誌輔助函數
package ginlogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogErrorWithStack 記錄錯誤日誌，包含 trace_id 和堆栈信息
func LogErrorWithStack(ctx *gin.Context, msg string, err error, extraFields ...zap.Field) {
	fields := make([]zap.Field, 0, len(extraFields)+2)

	// 添加 trace_id
	if traceID, exists := ctx.Get("trace_id"); exists {
		fields = append(fields, zap.String("trace_id", traceID.(string)))
	}

	// 添加錯誤
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, zap.Stack("stack")) // 添加堆栈信息

	// 添加額外字段
	fields = append(fields, extraFields...)

	// 記錄錯誤（不包含 Stack，因為 Console 模式下會壓縮成一行）
	zap.L().Error(msg, fields...)
}
