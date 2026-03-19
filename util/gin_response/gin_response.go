// Package ginresp 提供了通用的 Gin 回應結構和方法，方便在 Gin 框架中統一處理 API 回應。
package ginresp

import (
	"self_go_gin/common/msgid"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 通用回應結構
type Response struct {
	Result    msgid.MsgID `json:"result" binding:"required"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	TraceID   string      `json:"trace_id,omitempty"`  // 請求追蹤 ID
	Timestamp int64       `json:"timestamp,omitempty"` // 響應時間戳
}

// CreateMsgData 創建訊息資料
func CreateMsgData(key, value string) map[string]string {
	var msg = make(map[string]string)
	msg[key] = value
	return msg
}

// SuccessResponse 成功回應
func SuccessResponse(c *gin.Context, statusCode int, msg string, data interface{}, result msgid.MsgID) {
	traceID, _ := c.Get("trace_id")
	c.JSON(statusCode, Response{
		Result:    result,
		Msg:       msg,
		Data:      data,
		TraceID:   traceID.(string),
		Timestamp: time.Now().Unix(),
	})
}

// ErrorResponse 錯誤回應
func ErrorResponse(c *gin.Context, statusCode int, msg string, result msgid.MsgID, errData interface{}) {
	traceID, _ := c.Get("trace_id")
	c.JSON(statusCode, Response{
		Result:    result,
		Msg:       msg,
		Data:      errData,
		TraceID:   traceID.(string),
		Timestamp: time.Now().Unix(),
	})
}
