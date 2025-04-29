package mysql

import (
	"fmt"
	"gorm.io/gorm"
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

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 25
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSortBy() string {
	if p.SortBy == "" {
		p.SortBy = "id"
	}
	return p.SortBy
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "desc"
	}
	return p.Sort
}

func Paginate(pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(fmt.Sprintf("%s %s", pagination.GetSortBy(), pagination.GetSort()))
	}
}
