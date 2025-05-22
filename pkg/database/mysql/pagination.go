package mysql

import (
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"kaptan/pkg/utils"
	"kaptan/pkg/utils/dto"
	"math"
)

type Pagination struct {
	Limit       int    `json:"limit,omitempty;query:limit"`
	Page        int    `json:"page,omitempty;query:page"`
	Sort        string `json:"sort,omitempty;query:sort"`
	SortBy      string `json:"sort_by,omitempty;query:sort"`
	TotalRows   int64  `json:"total"`
	TotalPages  int    `json:"last_page"`
	CurrentPage int    `json:"current_page"`
	PerPage     int    `json:"per_page"`
}

func Paginate(pagination *Pagination, db *gorm.DB, paginationDto dto.Pagination) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = cast.ToInt(utils.If(totalPages < 1, 1, totalPages))
	pagination.Sort = paginationDto.GetSort()
	pagination.SortBy = paginationDto.GetSortBy()
	pagination.Limit = paginationDto.GetLimit()
	pagination.Page = paginationDto.GetPage()
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(paginationDto.GetOffset()).Limit(paginationDto.GetLimit()).Order(fmt.Sprintf("%s %s", paginationDto.GetSortBy(), paginationDto.GetSort()))
	}
}
