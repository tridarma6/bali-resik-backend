package helper

import (
	"math"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PaginationMeta struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type PaginationParam struct {
	Page    int `query:"page"`
	PerPage int `query:"per_page"`
}

func GetPaginationParam(ctx echo.Context) PaginationParam {
	page := 1
	perPage := 20

	p := new(PaginationParam)
	if err := ctx.Bind(p); err == nil {
		if p.Page > 0 {
			page = p.Page
		}
		if p.PerPage > 0 && p.PerPage <= 100 {
			perPage = p.PerPage
		}
	}

	return PaginationParam{
		Page:    page,
		PerPage: perPage,
	}
}

func Paginate(query *gorm.DB, param PaginationParam) (*gorm.DB, *PaginationMeta) {
	var total int64
	query.Count(&total)

	offset := (param.Page - 1) * param.PerPage

	totalPage := int(math.Ceil(float64(total) / float64(param.PerPage)))

	return query.Offset(offset).Limit(param.PerPage), &PaginationMeta{
		Page:      param.Page,
		PerPage:   param.PerPage,
		Total:     total,
		TotalPage: totalPage,
	}
}
