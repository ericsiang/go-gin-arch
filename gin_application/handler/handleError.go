// Package handler 處理錯誤
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"self_go_gin/common/msgid"
	ginlogger "self_go_gin/gin_application/inter/log"
	ginresp "self_go_gin/gin_application/inter/response"
	apperror "self_go_gin/internal/apperror"

	mysqlmgr "self_go_gin/util/mysql_manager"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandleError 統一錯誤處理 - 支持 AppError 和標準錯誤
func HandleError(context *gin.Context, err error) (bool, error) {
	if err == nil {
		return true, nil
	}

	// 檢查是否為 AppError
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return handleAppError(context, appErr)
	}

	// 處理標準錯誤
	return handleStandardError(context, err)
}

// handleAppError 處理業務應用錯誤
func handleAppError(context *gin.Context, appErr *apperror.AppError) (bool, error) {
	// 根據業務碼決定 HTTP 狀態碼和訊息
	httpStatus := getHTTPStatus(appErr.Code)
	msg := appErr.Message
	if msg == "" {
		msg = getMsgFromCode(appErr.Code)
	}

	// 構建詳細的錯誤日誌消息（包含 LogData）
	var logMsg string
	if len(appErr.LogData) > 0 {
		logMsg = fmt.Sprintf("AppError [Code: %v, Message: %s, LogData: %v]", appErr.Code, msg, appErr.LogData)
	} else {
		logMsg = fmt.Sprintf("AppError [Code: %v, Message: %s]", appErr.Code, msg)
	}

	// 記錄錯誤及堆棧
	ginlogger.LogErrorWithStack(context, logMsg, appErr.InternalErr)

	// 返回錯誤響應
	ginresp.ErrorResponse(context, httpStatus, msg, appErr.Code, nil)
	return false, fmt.Errorf("HandleError() AppError: %w", appErr)
}

// handleStandardError 處理標準錯誤（保留向後兼容性）
func handleStandardError(context *gin.Context, err error) (bool, error) {
	switch {
	case mysqlmgr.MysqlErrCode(err) == mysqlmgr.DuplicateEntryCode:
		ginresp.ErrorResponse(context, http.StatusBadRequest, "資源重複", msgid.DuplicateEntry, nil)
		return false, fmt.Errorf("HandleError() DuplicateEntryCode : %w", err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginresp.ErrorResponse(context, http.StatusNotFound, "資源不存在", msgid.NoContent, nil)
		return false, fmt.Errorf("HandleError() ErrRecordNotFound : %w", err)
	case errors.Is(err, ErrResourceExist):
		ginresp.ErrorResponse(context, http.StatusBadRequest, "資源已存在", msgid.ResourceExist, nil)
		return false, fmt.Errorf("HandleError() ErrResourceExist : %w", err)
	default:
		ginlogger.LogErrorWithStack(context, "Unexpected error", err)
		ginresp.ErrorResponse(context, http.StatusInternalServerError, "內部伺服器錯誤", msgid.Fail, nil)
		return false, fmt.Errorf("HandleError() : %w", err)
	}
}

// getHTTPStatus 根據業務碼返回 HTTP 狀態碼
func getHTTPStatus(code msgid.MsgID) int {
	switch code {
	case msgid.Success:
		return http.StatusOK
	case msgid.NoContent:
		return http.StatusNotFound
	case msgid.InvalidInput, msgid.DuplicateEntry, msgid.ResourceExist:
		return http.StatusBadRequest
	case msgid.TokenExpires, msgid.TokenInvalid:
		return http.StatusUnauthorized
	case msgid.RuleNotAllow:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// getMsgFromCode 根據業務碼返回默認訊息
func getMsgFromCode(code msgid.MsgID) string {
	switch code {
	case msgid.Success:
		return "操作成功"
	case msgid.Fail:
		return "操作失敗"
	case msgid.TokenExpires:
		return "令牌已過期"
	case msgid.TokenInvalid:
		return "令牌無效"
	case msgid.NoContent:
		return "資源不存在"
	case msgid.DuplicateEntry:
		return "資源已存在（重複條目）"
	case msgid.RuleNotAllow:
		return "不允許的操作"
	case msgid.ResourceExist:
		return "資源已存在"
	case msgid.InvalidInput:
		return "輸入驗證失敗"
	default:
		return "未知錯誤"
	}
}
