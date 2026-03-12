package gin_response

import (
	"self_go_gin/common/msgid"

	"github.com/gin-gonic/gin"
)

// Response 通用回應結構
type Response struct {
	Result msgid.MsgID `json:"result" binding:"required"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// CreateMsgData 創建訊息資料
func CreateMsgData(key, value string) map[string]string {
	var msg = make(map[string]string)
	msg[key] = value
	return msg
}

// SuccessResponse 成功回應
func SuccessResponse(c *gin.Context, statusCode int, msg string, data interface{}, result msgid.MsgID) {
	c.JSON(statusCode, Response{
		Result: result,
		Msg:    msg,
		Data:   data,
	})
}

// ErrorResponse 錯誤回應
func ErrorResponse(c *gin.Context, statusCode int, msg string, result msgid.MsgID, errData interface{}) {
	c.JSON(statusCode, Response{
		Result: result,
		Msg:    msg,
		Data:   errData,
	})
}
