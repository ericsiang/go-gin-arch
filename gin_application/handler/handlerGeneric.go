package handler

import (
	"fmt"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/common/appmsg"
	ginresp "self_go_gin/gin_application/inter/response"
	apperror "self_go_gin/internal/apperror"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	// ErrResourceNotFound 資源不存在
	ErrResourceNotFound = errors.New("resource_not_found")
	// ErrAccountLocked 帳號被鎖定
	ErrAccountLocked = errors.New("account_is_lock")
	// ErrDeleteNotAllow 刪除不允許
	ErrDeleteNotAllow = errors.New("delete_not_allow")
	// ErrResourceExist 資源已存在
	ErrResourceExist = errors.New("resource_exist")
	// ErrPasswordNoMatch 密碼不匹配
	ErrPasswordNoMatch = errors.New("password_not_match")
)

// GetHandler 處理獲取請求
func GetHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ginresp.FailResponse{})
		return false, apperror.NewAppError(
			msgid.Fail,
			appmsg.HandlerGetDataFailed,
			fmt.Errorf("GetHandler() error: %w", err),
			nil,
		)
	}
	return true, nil
}

// CreateHandler 處理創建請求
func CreateHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		mysqlErrorCheck, err := MysqlErrorCheck(ctx, err)
		if mysqlErrorCheck {
			return false, apperror.NewAppError(
				msgid.Fail,
				appmsg.HandlerMySQLError,
				fmt.Errorf("CreateHandler() MySQL error: %w", err),
				nil,
			)
		}
		ctx.JSON(http.StatusInternalServerError, ginresp.FailResponse{})
		return false, apperror.NewAppError(
			msgid.Fail,
			appmsg.HandlerCreateDataFailed,
			fmt.Errorf("CreateHandler() error: %w", err),
			nil,
		)

	}
	return true, nil
}

// UpdateHandler 處理更新請求
func UpdateHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, ginresp.FailResponse{
				Msg: "record_not_found",
			})
			return false, apperror.NewAppError(
				msgid.NoContent,
				appmsg.HandlerRecordNotFound,
				fmt.Errorf("UpdateHandler() record not found: %w", err),
				nil,
			)
		}
		mysqlErrorCheck, err := MysqlErrorCheck(ctx, err)
		if mysqlErrorCheck {
			return false, apperror.NewAppError(
				msgid.Fail,
				appmsg.HandlerMySQLError,
				fmt.Errorf("UpdateHandler() MySQL error: %w", err),
				nil,
			)
		}
		ctx.JSON(http.StatusInternalServerError, ginresp.FailResponse{})
		return false, apperror.NewAppError(
			msgid.Fail,
			appmsg.HandlerUpdateDataFailed,
			fmt.Errorf("UpdateHandler() error: %w", err),
			nil,
		)

	}
	return true, nil
}

// DeleteHandler 處理刪除請求
func DeleteHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		if errors.Is(err, ErrDeleteNotAllow) {
			ctx.JSON(http.StatusAccepted, ginresp.FailResponse{
				Msg: ErrDeleteNotAllow.Error(),
			})
			return false, apperror.NewAppError(
				msgid.Fail,
				appmsg.HandlerDeleteNotAllow,
				fmt.Errorf("DeleteHandler() delete not allowed: %w", err),
				nil,
			)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, ginresp.FailResponse{
				Msg: "record_not_found",
			})
			return false, apperror.NewAppError(
				msgid.NoContent,
				appmsg.HandlerRecordNotFound,
				fmt.Errorf("DeleteHandler() record not found: %w", err),
				nil,
			)
		} else if errors.Is(err, ErrResourceNotFound) {
			ctx.JSON(http.StatusNotFound, ginresp.FailResponse{
				Msg: ErrResourceNotFound.Error(),
			})
			return false, apperror.NewAppError(
				msgid.NoContent,
				appmsg.HandlerResourceNotFound,
				fmt.Errorf("DeleteHandler() resource not found: %w", err),
				nil,
			)
		}
		ctx.JSON(http.StatusInternalServerError, ginresp.FailResponse{})
		return false, apperror.NewAppError(
			msgid.Fail,
			appmsg.HandlerDeleteDataFailed,
			fmt.Errorf("DeleteHandler() error: %w", err),
			nil,
		)

	}
	return true, nil
}
