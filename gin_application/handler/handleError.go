package handler

import (
	"errors"
	"fmt"
	"net/http"
	"self_go_gin/common/msgid"
	"self_go_gin/util/gin_response"
	"self_go_gin/util/mysql_manager"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandlerError 處理錯誤
func HandlerError(context *gin.Context, err error) (bool, error) {
	switch {
	case mysql_manager.MysqlErrCode(err) == mysql_manager.DuplicateEntryCode:
		gin_response.ErrorResponse(context, http.StatusBadRequest, "", msgid.DuplicateEntry, nil)
		return false, fmt.Errorf("HandlerError() DuplicateEntryCode : %w", err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		gin_response.ErrorResponse(context, http.StatusNotFound, "", msgid.NoContent, nil)
		return false, fmt.Errorf("HandlerError() ErrRecordNotFound : %w", err)
	case errors.Is(err, ErrResourceExist):
		gin_response.ErrorResponse(context, http.StatusBadRequest, "", msgid.Fail, nil)
		return false, fmt.Errorf("HandlerError() ErrResourceExist : %w", err)
	case err != nil:
		gin_response.ErrorResponse(context, http.StatusInternalServerError, "", msgid.Fail, nil)
		return false, fmt.Errorf("HandlerError() : %w", err)
	default:
		return true, nil
	}
}
