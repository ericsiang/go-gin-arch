// Package apperror 錯誤處理結構和相關函數
package apperror

import (
	"fmt"
	"self_go_gin/common/msgid"

	"github.com/pkg/errors"
)

// AppError 定義錯誤標準結構
type AppError struct {
	Code        msgid.MsgID            // 業務代碼 (API 返回)
	Message     string                 // 給用戶看的訊息
	InternalErr error                  // 原始錯誤 (開發除錯用)
	LogData     map[string]interface{} // 額外的日誌數據 (用於調試)
}

// NewAppError 創建新的 AppError，並捕獲當前堆疊信息
func NewAppError(code msgid.MsgID, msg string, err error , logData map[string]interface{}) *AppError {
	return &AppError{
		Code:        code,
		Message:     msg,
		InternalErr: errors.WithStack(err),
		LogData:     logData,
	}
}

// NewAppErrorWithLogData 創建新的 AppError 並包含日誌數據
func NewAppErrorWithLogData(code msgid.MsgID, msg string, err error, logData map[string]interface{}) *AppError {
	return &AppError{
		Code:        code,
		Message:     msg,
		InternalErr: errors.WithStack(err),
		LogData:     logData,
	}
}

func (e *AppError) Error() string {
	if len(e.LogData) > 0 {
		return fmt.Sprintf("%s | LogData: %v", e.Message, e.LogData)
	}
	return e.Message
}
