package pagination

import (
	"fmt"
	"math"
	"strconv"
)

type PageInfo struct {
	Page        int `json:"page"`         // 页号
	PageSize    int `json:"page_size"`    // 每页大小
	TotalNumber int `json:"total_number"` // 数据总条数
	TotalPage   int `json:"total_page"`   // 数据总页数
}

func (p *PageInfo) SetTotalPage() *PageInfo {
	page, _ := strconv.Atoi(fmt.Sprintf("%1.0f", math.Ceil(float64(p.TotalNumber)/float64(p.PageSize))))
	p.TotalPage = page
	return p
}
