package pagination

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Paginator struct {
	DB         *gorm.DB
	Page       int
	PageSize   int
	OriginList interface{}
	TargetList interface{}
}

func PaginatorBuilder(builder Paginator) (pageInfo PageInfo, err error) {
	var total int64
	builder.DB.Count(&total)

	offset := (builder.Page - 1) * builder.PageSize
	builder.DB.Offset(offset).Limit(builder.PageSize).Find(builder.OriginList)
	if builder.DB.Error != nil {
		return pageInfo, builder.DB.Error
	}
	if builder.TargetList != nil {
		str, err := json.Marshal(builder.OriginList)
		if err != nil {
			return pageInfo, err
		}
		err = json.Unmarshal(str, &builder.TargetList)
		if err != nil {
			return pageInfo, err
		}
	}

	pageInfo = PageInfo{
		Page:        builder.Page,
		PageSize:    builder.PageSize,
		TotalNumber: int(total),
	}
	pageInfo.SetTotalPage()

	return
}
