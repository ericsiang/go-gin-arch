// Package ginlogger 提供 Zap 日誌輔助函數
package ginlogger

import (
	"self_go_gin/internal/apperror"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogErrorWithStack 記錄錯誤日誌，包含 trace_id 和堆栈信息
func LogErrorWithStack(ctx *gin.Context, msg string, appErr *apperror.AppError, extraFields ...zap.Field) {

	// 🔥 輸出完整的錯誤鏈到日誌
	chain := appErr.ErrorChain()

	// 預分配足夠的容量：trace_id + error + stack + extraFields + AppError 字段
	fields := make([]zap.Field, 0, len(extraFields)+7)

	// 添加 trace_id
	if traceID, exists := ctx.Get("trace_id"); exists {
		fields = append(fields, zap.String("trace_id", traceID.(string)))
	}

	// 添加 AppError 特定字段到 fields
	fields = append(fields, zap.Uint32("code", uint32(appErr.Code)))
	fields = append(fields, zap.String("message", appErr.Message))
	fields = append(fields, zap.String("layer", appErr.Layer))
	fields = append(fields, zap.Any("log_data", appErr.LogData))
	fields = append(fields, zap.Strings("chain", chain))

	// 添加堆棧信息
	fields = append(fields, zap.Stack("stack"))

	// 添加額外字段
	fields = append(fields, extraFields...)

	// 記錄完整的錯誤日誌
	zap.L().Error(msg, fields...)

}
