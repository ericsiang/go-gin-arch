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
	InternalErr error                  // 原始錯誤 (開發除錯用) - 自動追踪調用棧
	LogData     map[string]interface{} // 額外的日誌數據 (用於調試)
	Layer       string                 // 錯誤發生的層級
}

// Option 是用於設置 AppError 選項的函數類型
type Option func(*AppError)

// NewAppError 創建新的 AppError
func NewAppError(code msgid.MsgID, msg string, err error, opts ...Option) *AppError {
	// 如果 err 已是 AppError，直接使用；否則用 errors.WithStack 追踪調用棧
	internalErr := err
	if _, ok := err.(*AppError); !ok && err != nil {
		internalErr = errors.WithStack(err) // ✨ 自動追踪（僅非 AppError 錯誤）
	}

	appErr := &AppError{
		Code:        code,
		Message:     msg,
		InternalErr: internalErr,
		LogData:     make(map[string]interface{}),
	}

	// 應用所有選項
	for _, opt := range opts {
		opt(appErr)
	}

	return appErr
}

// WithLayer 設置錯誤發生的層級
func WithLayer(layer string) Option {
	return func(e *AppError) {
		e.Layer = layer
	}
}

// WithLogData 設置日誌數據
func WithLogData(logData map[string]interface{}) Option {
	return func(e *AppError) {
		for k, v := range logData {
			e.LogData[k] = v
		}
	}
}

// Wrap 包裝錯誤為 AppError（避免重複包裝）
func Wrap(err error, code msgid.MsgID, msg string, layer string) *AppError {
	if err == nil {
		return nil
	}
	// 如果已經是 AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(code, msg, err, WithLayer(layer))
}

func (e *AppError) Error() string {
	if len(e.LogData) > 0 {
		return fmt.Sprintf("%s | LogData: %v", e.Message, e.LogData)
	}
	return e.Message
}

// Unwrap 返回底層錯誤（支持 errors.Is）
func (e *AppError) Unwrap() error {
	return e.InternalErr
}

// ErrorChain 返回完整的錯誤鏈（用於日誌輸出）
func (e *AppError) ErrorChain() []string {
	var chain []string

	// 添加當前 AppError
	if e.Layer != "" {
		chain = append(chain, fmt.Sprintf("[Code:%d] %s (Layer:%s)", e.Code, e.Message, e.Layer))
	} else {
		chain = append(chain, fmt.Sprintf("[Code:%d] %s", e.Code, e.Message))
	}

	// 遍歷包裝的錯誤鏈
	err := e.InternalErr
	for err != nil {
		if appErr, ok := err.(*AppError); ok {
			// 遞迴查看下層的 AppError
			if appErr.Layer != "" {
				chain = append(chain, fmt.Sprintf("=> [Code:%d] %s (Layer:%s)", appErr.Code, appErr.Message, appErr.Layer))
			} else {
				chain = append(chain, fmt.Sprintf("=> [Code:%d] %s", appErr.Code, appErr.Message))
			}
			err = appErr.InternalErr
		} else {
			// 標準錯誤，直接添加
			chain = append(chain, fmt.Sprintf("=> %v", err))
			break
		}
	}

	return chain
}
