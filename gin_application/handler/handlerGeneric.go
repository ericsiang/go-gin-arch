package handler

import (
	"fmt"

	"net/http"
	"self_go_gin/gin_application/inter/response"

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
		ctx.JSON(http.StatusInternalServerError, response.FailResponse{})
		return false, fmt.Errorf("GetHandler() \n %w", err)
	}
	return true, nil
}

// CreateHandler 處理創建請求
func CreateHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		mysqlErrorCheck, err := MysqlErrorCheck(ctx, err)
		if mysqlErrorCheck {
			return false, fmt.Errorf("CreateHandler() \n %w", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, response.FailResponse{})
			return false, fmt.Errorf("CreateHandler() \n %w", err)
		}
	}
	return true, nil
}

// UpdateHandler 處理更新請求
func UpdateHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.FailResponse{
				Msg: "record_not_found",
			})
			return false, fmt.Errorf("UpdateHandler() \n %w", err)
		}
		mysqlErrorCheck, err := MysqlErrorCheck(ctx, err)
		if mysqlErrorCheck {
			return false, fmt.Errorf("UpdateHandler() \n %w", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, response.FailResponse{})
			return false, fmt.Errorf("UpdateHandler() \n %w", err)
		}
	}
	return true, nil
}

// DeleteHandler 處理刪除請求
func DeleteHandler(ctx *gin.Context, err error) (bool, error) {
	if err != nil {
		if errors.Is(err, ErrDeleteNotAllow) {
			ctx.JSON(http.StatusAccepted, response.FailResponse{
				Msg: ErrDeleteNotAllow.Error(),
			})
			return false, fmt.Errorf("DeleteHandler() \n %w", err)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.FailResponse{
				Msg: "record_not_found",
			})
			return false, fmt.Errorf("DeleteHandler() \n %w", err)
		} else if errors.Is(err, ErrResourceNotFound) {
			ctx.JSON(http.StatusNotFound, response.FailResponse{
				Msg: ErrResourceNotFound.Error(),
			})
			return false, fmt.Errorf("DeleteHandler() \n %w", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, response.FailResponse{})
			return false, fmt.Errorf("DeleteHandler() \n %w", err)
		}
	}
	return true, nil
}
