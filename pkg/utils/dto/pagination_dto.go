package dto

type Pagination struct {
	Page       int64 `json:"page" form:"page"`
	Limit      int64 `json:"limit" form:"limit"`
	Pagination bool  `json:"is_paginated" form:"is_paginated"`
}

func (p *Pagination) SetDefault() {
	p.Limit = 25
	p.Page = 1
	p.Pagination = true
}
