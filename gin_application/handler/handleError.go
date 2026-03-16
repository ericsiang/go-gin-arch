// Package handler 處理錯誤
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

// HandleError 處理錯誤
func HandleError(context *gin.Context, err error) (bool, error) {
	switch {
	case mysqlmgr.MysqlErrCode(err) == mysqlmgr.DuplicateEntryCode:
		ginresp.ErrorResponse(context, http.StatusBadRequest, "", msgid.DuplicateEntry, nil)
		return false, fmt.Errorf("HandlerError() DuplicateEntryCode : %w", err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginresp.ErrorResponse(context, http.StatusNotFound, "", msgid.NoContent, nil)
		return false, fmt.Errorf("HandlerError() ErrRecordNotFound : %w", err)
	case errors.Is(err, ErrResourceExist):
		ginresp.ErrorResponse(context, http.StatusBadRequest, "", msgid.ResourceExist, nil)
		return false, fmt.Errorf("HandlerError() ErrResourceExist : %w", err)
	case err != nil:
		ginresp.ErrorResponse(context, http.StatusInternalServerError, "", msgid.Fail, nil)
		return false, fmt.Errorf("HandlerError() : %w", err)
	default:
		return true, nil
	}
}
