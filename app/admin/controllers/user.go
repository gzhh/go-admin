package controllers

import (
	"fmt"
	"go-admin/app/admin/handlers"
	"go-admin/app/admin/validators"
	"go-admin/pkg/service/admin"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserLogin
// @Summary 用户-登录
// @Description 返回结构中有token，则说明授权成功；若无token，则说明账号在审核。
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param body body validators.UserLogin true "数据参数"
// @Success 200 {object} handlers.resp{data=admin.UserResponse}
// @Router /user/login [post]
func UserLogin(ctx *gin.Context) {
	// validate params
	var params validators.UserLogin
	if err := ctx.ShouldBindJSON(&params); err != nil {
		handlers.Fail(ctx, err.Error())
		return
	}

	if err := validators.CheckCtx(ctx, params); err != nil {
		handlers.Fail(ctx, err.Error())
		return
	}

	srv := admin.NewUser()
	data, err := srv.UserLogin(ctx.Request.Context(), params)
	if err != nil {
		handlers.Fail(ctx, fmt.Sprintf("登录失败, %v", err))
		return
	}

	// response
	handlers.Success(ctx, data)
}

// UserList
// @Summary 用户-列表接口
// @Description
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param raw query validators.UserList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} handlers.resp{data=admin.UserListResponse}
// @Router /user [get]
func UserList(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Request.Header.Get("user_id"))

	// validate params
	var params validators.UserList
	if err := ctx.ShouldBind(&params); err != nil {
		handlers.Fail(ctx, err.Error())
		return
	}

	if err := validators.CheckCtx(ctx, params); err != nil {
		handlers.Fail(ctx, err.Error())
		return
	}

	srv := admin.NewUser()
	ok, err := srv.IsAdmin(ctx.Request.Context(), userID)
	if err != nil {
		handlers.Fail(ctx, err.Error())
		return
	}
	if !ok {
		handlers.Fail(ctx, "You are not admin")
		return
	}
	// handle logic
	data, err := srv.GetList(ctx.Request.Context(), params)
	if err != nil {
		handlers.Fail(ctx, fmt.Sprintf("列表获取失败, %s", err.Error()))
		return
	}

	// response
	handlers.Success(ctx, data)
}
