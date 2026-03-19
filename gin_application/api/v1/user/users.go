// Package v1 用戶相關API
package v1

import (
	"fmt"
	"net/http"

	"self_go_gin/common/msgid"
	"self_go_gin/domains/user/service"
	"self_go_gin/gin_application/api/v1/user/request"
	"self_go_gin/gin_application/handler"
	ginlogger "self_go_gin/gin_application/inter/log"
	ginresp "self_go_gin/util/gin_response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CreateUser 創建用戶
// @Summary  Create Users
// @Description Create Users
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body swagger_docs.DocUsersCreate true "request body"
// @Success 200 {string} json "{"msg": {"success": "success"},"data": {}}"
// @Failure 400 {string} json "{"msg": {"fail": "帳密錯誤"},"data": null}"
// @Router /api/v1/auth/users [post]
func CreateUser(ctx *gin.Context) {
	var data request.CreateUserRequest
	// var respData response.CreateUserResponse

	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			// gin_response.ErrorResponse(ctx, http.StatusBadRequest, "request_parameter_validation_failed", common_msg_id.Fail, nil)
			return
		}
		// 非validator.ValidationErrors類型錯誤直接傳回
		ginlogger.LogErrorWithStack(ctx, "Api CreateUser() ShouldBindBodyWith fail", err)
		ginresp.ErrorResponse(ctx, http.StatusNotFound, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	userService, err := service.NewUserService()
	if err != nil {
		ginlogger.LogErrorWithStack(ctx, "Api CreateUser() NewUserService fail", err)
		ginresp.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	_, err = userService.CreateUser(data)
	ok, err := handler.HandleError(ctx, err)
	if !ok {
		ginlogger.LogErrorWithStack(ctx, "Api CreateUser() CreateUser fail", err)
		return	
	}
	ginresp.SuccessResponse(ctx, http.StatusOK, "", nil, msgid.Success)
}

// UserLogin 用戶登入
// @Summary  User Login
// @Description User Login
// @Tags Users
// @Accept json
// @Produce json
// @Param request body swagger_docs.DocUsersLogin true "request body"
// @Success 200 {string}  "成功"
// @Failure 400 {string}  "失敗"
// @Failure 401 {string}  "Unauthorized"
// @Router /api/v1/users/login [post]
func UserLogin(ctx *gin.Context) {
	var data request.UserLoginRequest
	if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		check := handler.ValidCheckAndTrans(ctx, err)
		if check {
			ginresp.ErrorResponse(ctx, http.StatusBadRequest, "request_parameter_validation_failed", msgid.Fail, nil)
			return
		}
		// 非validator.ValidationErrors類型錯誤直接傳回
		ginlogger.LogErrorWithStack(ctx, "Api UserLogin() ShouldBindBodyWith fail", err)
		ginresp.ErrorResponse(ctx, http.StatusNotFound, "invalid_request_parameters", msgid.Fail, nil)
		return
	}

	userService, err := service.NewUserService()
	if err != nil {
		ginlogger.LogErrorWithStack(ctx, "Api UserLogin() NewUserService fail", err)
		ginresp.ErrorResponse(ctx, http.StatusInternalServerError, "internal_server_error", msgid.Fail, nil)
		return
	}
	jwtToken, err := userService.CheckLogin(data)
	ok, err := handler.HandleError(ctx, err)
	if !ok {
		ginlogger.LogErrorWithStack(ctx, "Api UserLogin()", err)
		return
	}
	ginresp.SuccessResponse(ctx, http.StatusOK, "", ginresp.CreateMsgData("jwt_token", *jwtToken), msgid.Success)
}

// GetUsersByID 根據ID獲取用戶
// @Summary Get Users By ID
// @Description Get Users By ID
// @Tags Users
// @Accept json
// @Produce json
// @Security JwtTokenAuth
// @Param filterUsersId path string true "filterUsersId"
// @Success 200 {string}  "成功"
// @Failure 400 {string}  "失敗"
// @Failure 401 {string}  "Unauthorized"
// @Router /api/v1/auth/users/{filterUsersId} [get]
func GetUsersByID(ctx *gin.Context) {
	var data request.GetUsersByIDRequest
	usersID, ok := ctx.Get("usersID")
	if !ok {
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "can not get users", msgid.Fail, nil)
		return
	}
	data.FilterUsersID = ctx.Param("filterUsersID")
	stringUsersID := fmt.Sprintf("%v", usersID)
	if data.FilterUsersID != stringUsersID {
		ginresp.ErrorResponse(ctx, http.StatusBadRequest, "user not match", msgid.Fail, nil)
		return
	}

	ginresp.SuccessResponse(ctx, http.StatusOK, "success", data.FilterUsersID, msgid.Success)

}
