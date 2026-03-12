package handler

import (
	"net/http"
	"self_go_gin/common/msgid"
	validlang "self_go_gin/gin_application/validate_lang"

	"self_go_gin/util/gin_response"

	"github.com/gin-gonic/gin"
)

// ValidCheckAndTrans 檢查驗證錯誤並翻譯
func ValidCheckAndTrans(context *gin.Context, err error) (ok bool) {
	validTrans, ok := validlang.ErrorValidateCheckAndTrans(err)
	if ok {
		var validTransMsgData = make(map[int]string)
		var i int
		for _, v := range validTrans {
			validTransMsgData[i] = v
			i++
		}

		gin_response.ErrorResponse(context, http.StatusBadRequest, "validate fail", msgid.Fail, validTransMsgData)
		return ok
	}
	return false
}
