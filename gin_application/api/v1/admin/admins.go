// Package v1 管理員相關API
package v1

import (
	"fmt"
	"net/http"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/admin/service"
	"self_go_gin/gin_application/api/v1/admin/request"
	"self_go_gin/gin_application/api/v1/admin/response"
	"self_go_gin/gin_application/handler"
	ginlogger "self_go_gin/gin_application/inter/log"
	ginresp "self_go_gin/gin_application/inter/response"
	"self_go_gin/internal/apperror"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "self_go_gin/gin_application/swaggerdocs"

)

// CreateAdmin 創建管理員
// @Summary  Create Admins
// @Description Create Admins
// @Tags Admins
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body swaggerdocs.DocAdminsCreate true "request body"
// @Success 200 {string} json "{"msg": {"success": "success"},"data": {}}"
// @Failure 400 {string} json "{"msg": {"fail": "帳密錯誤"},"data": null}"
// @Router /api/v1/admins [post]
func CreateAdmin(ctx *gin.Context) {
	var data request.CreateAdminRequest
	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			return
		}
		appError := apperror.NewAppError(msgid.NotAppError, "類型錯誤", err, nil)
		// 非validator.ValidationErrors類型錯誤直接傳回
		ginlogger.LogErrorWithStack(ctx, "Api CreateUser() ShouldBindBodyWith fail", appError)
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	adminService, err := service.NewAdminService()
	if err != nil {
		ginlogger.LogErrorWithStack(ctx, "Api CreateAdmin() NewAdminService() fail", err)
		ginresp.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	admin, err := adminService.CreateAdmin(data)
	ok, err := handler.HandleError(ctx, err)
	if !ok {
		ginlogger.LogErrorWithStack(ctx, "Api CreateAdmin() CreateAdmin() fail", err)
		return
	}

	respData := response.CreateAdminResponse{
		AdminID: admin.ID,
		Account: admin.GetAccount(),
	}
	ginresp.SuccessResponse(ctx, http.StatusOK, "", respData, msgid.Success)
}

// AdminLogin 管理員登入
// @Summary  Admin Login
// @Description Admin Login
// @Tags Admins
// @Accept json
// @Produce json
// @Param request body swaggerdocs.DocAdminsLogin true "request body"
// @Success 200 {string}  "成功"
// @Failure 400 {string}  "失敗"
// @Failure 401 {string}  "Unauthorized"
// @Router /api/v1/admins/login [post]
func AdminLogin(ctx *gin.Context) {
	var data request.AdminLoginRequest

	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			return
		}
		appError := apperror.NewAppError(msgid.NotAppError, "類型錯誤", err, nil)
		// 非validator.ValidationErrors類型錯誤直接傳回
		ginlogger.LogErrorWithStack(ctx, "Api AdminLogin() ShouldBindBodyWith fail", appError)
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	adminService, err := service.NewAdminService()
	if err != nil {
		ginlogger.LogErrorWithStack(ctx, "Api AdminLogin() NewAdminService fail", err)
		ginresp.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	jwtToken, err := adminService.CheckLogin(data)
	ok, err := handler.HandleError(ctx, err)
	if !ok {
		ginlogger.LogErrorWithStack(ctx, "Api AdminLogin() CheckLogin fail", err)
		return
	}
	ginresp.SuccessResponse(ctx, http.StatusOK, "", ginresp.CreateMsgData("jwt_token", *jwtToken), msgid.Success)

}

// GetAdminsByID 根據ID獲取管理員
// @Summary Get Admins By ID
// @Description Get Admins By ID
// @Tags Admins
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param filterAdminsId path string true "filterAdminsId"
// @Success 200 {string}  "成功"
// @Failure 400 {string}  "失敗"
// @Failure 401 {string}  "Unauthorized"
// @Router /api/v1/admins/{filterAdminsId} [get]
func GetAdminsByID(ctx *gin.Context) {
	var data request.GetAdminsByIDRequest

	adminID, ok := ctx.Get("adminID")
	if !ok {
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "can not get admins", msgid.Fail, nil)
		return
	}
	data.FilterAdminsID = ctx.Param("filterAdminsID")
	stringAdminsID := fmt.Sprintf("%v", adminID)
	if data.FilterAdminsID != stringAdminsID {
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "admin not match", msgid.Fail, nil)
		return
	}

	ginresp.SuccessResponse(ctx, http.StatusOK, "", nil, msgid.Success)
}
