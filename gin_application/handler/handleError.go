// Package handler 處理錯誤
package handler

import (
	"fmt"
	"net/http"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/common/appmsg"
	ginresp "self_go_gin/gin_application/inter/response"
	apperror "self_go_gin/internal/apperror"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	// ErrRecordNotFound 記錄未找到
	ErrRecordNotFound = errors.New("record not found")
	// ErrResourceExist 資源已存在
	ErrResourceExist = errors.New("resource exist")
	// ErrDuplicateEntry 重複條目
	ErrDuplicateEntry = errors.New("duplicate entry")
	// ErrPasswordIncorrect 密碼錯誤
	ErrPasswordIncorrect = errors.New("password incorrect")
	// ErrGetDBFailed 獲取資料庫實例失敗
	ErrGetDBFailed = errors.New("failed to get gorm.DB instance")
)

// HandleError 處理錯誤
func HandleError(ctx *gin.Context, err error) (bool, *apperror.AppError) {
	fmt.Println("HandleError err:", err)

	if err == nil {
		fmt.Println("no error")
		return true, nil
	}
	err = MysqlErrorCheck(err)

	if appErr, ok := err.(*apperror.AppError); ok {
		// 根據業務代碼返回適當的 HTTP 狀態碼
		httpStatus := getHTTPStatus(appErr.Code)

		// 返回給客戶端
		ginresp.ErrorResponse(ctx, httpStatus, appErr.Message, appErr.Code, appErr.LogData)
		return false, appErr
	}

	return false, apperror.NewAppError(msgid.NotAppError, "錯誤包裝", err, nil)

	// 處理其他特定錯誤類型
	// switch {
	// case mysqlmgr.MysqlErrCode(err) == mysqlmgr.DuplicateEntryCode:
	// 	ginresp.ErrorResponse(context, http.StatusBadRequest, "資源重複", msgid.DuplicateEntry, nil)
	// 	return false, fmt.Errorf("HandlerError() DuplicateEntryCode : %w", err)
	// case errors.Is(err, gorm.ErrRecordNotFound):
	// 	ginresp.ErrorResponse(context, http.StatusNotFound, "資源不存在", msgid.NoContent, nil)
	// 	return false, fmt.Errorf("HandlerError() ErrRecordNotFound : %w", err)
	// case errors.Is(err, ErrResourceExist):
	// 	ginresp.ErrorResponse(context, http.StatusBadRequest, "資源已存在", msgid.ResourceExist, nil)
	// 	return false, fmt.Errorf("HandlerError() ErrResourceExist : %w", err)
	// default:
	// 	ginresp.ErrorResponse(context, http.StatusInternalServerError, "", msgid.Fail, nil)
	// 	return false, fmt.Errorf("HandlerError() : %w", err)
	// }
}

// MysqlErrorCheck 檢查 MySQL 錯誤
func MysqlErrorCheck(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 { // Duplicate entry detected
			return apperror.NewAppError(
				msgid.DuplicateEntry,
				appmsg.DuplicateEntry,
				err,
				apperror.WithLayer("Handler MysqlErrorCheck()"),
			)
		}
	}
	return err
}

// getHTTPStatus 根據業務代碼對應到 HTTP 狀態碼
func getHTTPStatus(code msgid.MsgID) int {
	switch code {
	case msgid.InvalidInput, msgid.DuplicateEntry:
		return http.StatusBadRequest
	case msgid.NoContent:
		return http.StatusNotFound
	case msgid.ResourceExist:
		return http.StatusConflict
	case msgid.TokenExpires, msgid.TokenInvalid:
		return http.StatusUnauthorized
	case msgid.RuleNotAllow:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
