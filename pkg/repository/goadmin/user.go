package goadmin

import (
	"go-admin/pkg/repository"
	"time"
)

const (
	// role_type
	AdminUserRoleTypeNormal = 0
	AdminUserRoleTypeAdmin  = 1

	// status
	AdminUserStatusDisable = 0
	AdminUserStatusEnable  = 1
)

var AdminUserRoleTypeLabel = []map[string]interface{}{
	{
		"id":   AdminUserRoleTypeNormal,
		"name": "普通用户",
	},
	{
		"id":   AdminUserRoleTypeAdmin,
		"name": "管理员",
	},
}

const TableNameAdminUser = "admin_user"

// AdminUser mapped from table <admin_user>
type AdminUser struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username   string    `gorm:"column:username;not null;comment:用户名" json:"username"`                                  // 用户名
	Password   string    `gorm:"column:password;not null;comment:密码" json:"password"`                                   // 密码
	Sub        string    `gorm:"column:sub;not null;comment:第三方登录sub" json:"sub"`                                       // 第三方登录sub
	Remark     string    `gorm:"column:remark;not null;comment:备注" json:"remark"`                                       // 备注
	Status     int32     `gorm:"column:status;not null;comment:使用状态：0-停用 1-启用" json:"status"`                           // 使用状态：0->
	LoginType  int32     `gorm:"column:login_type;not null;comment:登录类型：0账号密码 1钉钉" json:"login_type"`                   // 登录类型：0
	RoleType   int32     `gorm:"column:role_type;not null;comment:角色：0普通用户 1管理员" json:"role_type"`                      // 角色：0普通>
	IsDelete   int32     `gorm:"column:is_delete;not null;comment:是否删除：0否 1是" json:"is_delete"`                         // 是否删除：0否 1
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"` // 更新时间
}

// TableName AdminUser's table name
func (*AdminUser) TableName() string {
	return TableNameAdminUser
}

func NewAdminUser() *AdminUserModel {
	repo := &AdminUserModel{}
	repo.DriverKey = DriverKeyGoAdmin
	repo.Table = TableNameAdminUser
	return repo
}

type AdminUserModel struct {
	repository.Base
}

func (repo *AdminUserModel) GetByUsername(username string) (AdminUser, error) {
	var user AdminUser
	db, err := repo.GetDBSession()
	if err != nil {
		return user, err
	}
	db.Select("*")
	db.Where("is_delete = ? and username = ?", IsDeleteFalse, username).First(&user)
	return user, db.Error
}
