package validators

type UserLogin struct {
	LoginType *int   `json:"login_type" validate:"required"` // 注册类型：0-账号密码 1-钉钉
	Username  string `json:"username" validate:"omitempty"`  // 用户名
	Password  string `json:"password" validate:"omitempty"`  // 密码
	Code      string `json:"code" validate:"omitempty"`      // 第三方登录临时授权码code
}

type UserList struct {
	Page     int `form:"page" validate:"required,min=1"`               // 页号
	PageSize int `form:"page_size" validate:"required,min=1,max=1000"` // 每页大小
}
