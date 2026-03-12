package v1

import (
	"fmt"
	"net/http"
	"self_go_gin/common/msgid"
	"self_go_gin/domains/admin/service"
	"self_go_gin/gin_application/api/v1/admin/request"
	"self_go_gin/gin_application/api/v1/admin/response"
	"self_go_gin/gin_application/handler"
	"self_go_gin/util/gin_response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// CreateAdmin 創建管理員
// @Summary  Create Admins
// @Description Create Admins
// @Tags Admins
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body swagger_docs.DocAdminsCreate true "request body"
// @Success 200 {string} json "{"msg": {"success": "success"},"data": {}}"
// @Failure 400 {string} json "{"msg": {"fail": "帳密錯誤"},"data": null}"
// @Router /api/v1/auth/admins [post]
func CreateAdmin(ctx *gin.Context) {
	var data request.CreateAdminRequest
	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			gin_response.ErrorResponse(ctx, http.StatusBadRequest, "request_parameter_validation_failed", msgid.Fail, nil)
			return
		}
		// 非validator.ValidationErrors類型錯誤直接傳回
		zap.L().Error("\n Api CreateAdmin() 失敗(ShouldBindBodyWith fail) : " + err.Error())
		gin_response.ErrorResponse(ctx, http.StatusNotFound, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	adminService ,err := service.NewAdminService()
	if err != nil {
		zap.L().Error("\n Api CreateAdmin() NewAdminService fail : " + err.Error())
		gin_response.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	admin, err := adminService.CreateAdmin(data)
	ok, err := handler.HandlerError(ctx, err)
	if !ok {
		zap.L().Error("\n Api CreateAdmin() \n " + err.Error())
		return
	}

	respData := response.CreateAdminResponse{
		AdminID: admin.ID,
		Account: admin.Account,
	}
	gin_response.SuccessResponse(ctx, http.StatusOK, "", respData, msgid.Success)
}

// AdminLogin 管理員登入
// @Summary  Admin Login
// @Description Admin Login
// @Tags Admins
// @Accept json
// @Produce json
// @Param request body swagger_docs.DocAdminsLogin true "request body"
// @Success 200 {string}  "成功"
// @Failure 400 {string}  "失敗"
// @Failure 401 {string}  "Unauthorized"
// @Router /api/v1/admins/login [post]
func AdminLogin(ctx *gin.Context) {
	var data request.AdminLoginRequest

	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			gin_response.ErrorResponse(ctx, http.StatusBadRequest, "request_parameter_validation_failed", msgid.Fail, nil)
			return
		}
		// 非validator.ValidationErrors類型錯誤直接傳回
		zap.L().Error("\n Api AdminLogin() 失敗(ShouldBindBodyWith fail) : " + err.Error())
		gin_response.ErrorResponse(ctx, http.StatusNotFound, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	adminService ,err := service.NewAdminService()
	if err != nil {
		zap.L().Error("\n Api AdminLogin() NewAdminService fail : " + err.Error())
		gin_response.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	jwtToken, err := adminService.CheckLogin(data)
	ok, err := handler.HandlerError(ctx, err)
	if !ok {
		zap.L().Error("\n Api AdminLogin() \n " + err.Error())
		return
	}
	gin_response.SuccessResponse(ctx, http.StatusOK, "", gin_response.CreateMsgData("jwt_token", *jwtToken), msgid.Success)

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
// @Router /api/v1/auth/admins/{filterAdminsId} [get]
func GetAdminsByID(ctx *gin.Context) {
	var data request.GetAdminsByIDRequest

	adminID, ok := ctx.Get("adminID")
	if !ok {
		gin_response.ErrorResponse(ctx, http.StatusBadRequest, "can not get admins", msgid.Fail, nil)
		return
	}
	data.FilterAdminsID = ctx.Param("filterAdminsID")
	stringAdminsID := fmt.Sprintf("%v", adminID)
	if data.FilterAdminsID != stringAdminsID {
		gin_response.ErrorResponse(ctx, http.StatusBadRequest, "admin not match", msgid.Fail, nil)
		return
	}

	gin_response.SuccessResponse(ctx, http.StatusOK, "", nil, msgid.Success)
}
