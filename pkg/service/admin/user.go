package admin

import (
	"context"
	"errors"
	"fmt"
	"go-admin/app/admin/validators"
	"go-admin/internal/lib/config"
	"go-admin/pkg/repository/goadmin"
	"go-admin/pkg/utils/auth"
	"go-admin/pkg/utils/pagination"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserLoginTypePassword = 0 // 账号密码登录
)

type user struct {
	adminUserRepo *goadmin.AdminUserModel
}

func NewUser() *user {
	return &user{
		adminUserRepo: goadmin.NewAdminUser(),
	}
}

type UserListResponse struct {
	List     []UserResponse      `json:"list"`      // 数据列表
	PageInfo pagination.PageInfo `json:"page_info"` // 分页信息
}

type UserResponse struct {
	ID         int32     `json:"id"`          // ID
	Username   string    `json:"username"`    // 用户名
	Remark     string    `json:"remark"`      // 备注
	Status     int32     `json:"status"`      // 角色：0停用 1启用
	RoleType   int32     `json:"role_type"`   // 角色：0普通用户 1管理员
	IsDelete   int32     `json:"is_delete"`   // 是否删除：0否 1是
	Token      string    `json:"token"`       // Token
	CreateTime time.Time `json:"create_time"` // 创建时间
}

func (srv user) GetList(ctx context.Context, params validators.UserList) (data UserListResponse, err error) {
	db, err := srv.adminUserRepo.GetDBSession()
	if err != nil {
		return UserListResponse{}, err
	}

	db.Select("*")
	db = db.Where("is_delete = ?", goadmin.IsDeleteFalse)
	db = db.Order("id DESC")

	var originList []goadmin.AdminUser
	var list []UserResponse
	pageInfo, err := pagination.PaginatorBuilder(pagination.Paginator{
		DB:         db,
		Page:       params.Page,
		PageSize:   params.PageSize,
		OriginList: &originList,
		TargetList: &list,
	})
	if err != nil {
		return UserListResponse{}, err
	}

	data.List = list
	data.PageInfo = pageInfo
	return data, nil
}

func (srv user) UserLogin(ctx context.Context, params validators.UserLogin) (resp UserResponse, err error) {
	if *params.LoginType == UserLoginTypePassword {
		// verify password
		var valid bool
		valid, err = srv.VerifyPassword(ctx, params)
		if err != nil {
			return resp, errors.New("invalid admin username/password")
		}
		if !valid {
			return resp, errors.New("error admin username/password")
		}

		// get token
		resp, err = srv.GetToken(ctx, params.Username)
		if err != nil {
			return resp, fmt.Errorf("generate token error, %v", err)
		}
	} else {
		return resp, errors.New("unsupported type")
	}

	return resp, nil
}

func (srv user) VerifyPassword(ctx context.Context, params validators.UserLogin) (valid bool, err error) {
	// get hashPassword from db
	cond := map[string]interface{}{
		"username":  params.Username,
		"is_delete": goadmin.IsDeleteFalse,
	}
	record, err := srv.adminUserRepo.Get("password", cond)
	if err != nil {
		return false, err
	}
	hashPassword := record["password"].(string)

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(params.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (srv user) GetToken(ctx context.Context, username string) (resp UserResponse, err error) {
	// get one
	userModel, err := srv.adminUserRepo.GetByUsername(username)
	if err != nil {
		return resp, fmt.Errorf("not found admin user[%s], error: %v", username, err.Error())
	}
	userInfo := map[string]interface{}{
		"id":        userModel.ID,
		"username":  userModel.Username,
		"remark":    userModel.Remark,
		"role_type": userModel.RoleType,
	}
	resp = UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Remark:   userModel.Remark,
		RoleType: userModel.RoleType,
	}

	// generate token
	token, err := auth.GenerateToken(config.Settings.AdminServer.Server.JwtSecret, userInfo)
	if err != nil {
		return resp, fmt.Errorf("generate token error: %v", err.Error())
	}
	// userInfo["token"] = token
	resp.Token = token

	return resp, nil
}

func (srv user) IsAdmin(ctx context.Context, userID int) (bool, error) {
	dbSession, err := srv.adminUserRepo.GetDBSession()
	if err != nil {
		return false, err
	}

	condition := map[string]interface{}{
		"id":        userID,
		"role_type": goadmin.AdminUserRoleTypeAdmin,
	}

	var exist bool
	err = dbSession.Where(condition).Select("count(*) > 0").Find(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, nil
}
