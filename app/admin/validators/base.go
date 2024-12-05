package validators

type CommonList struct {
	Page      int `form:"page" validate:"required,min=1"`               // 页号
	PageSize  int `form:"page_size" validate:"required,min=1,max=1000"` // 页大小
	QueryTime int `form:"query_time" validate:"omitempty"`              // 查询截止时间戳
}
