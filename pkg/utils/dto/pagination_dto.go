package dto

type Pagination struct {
	Page       int64 `json:"page" form:"page"`
	Limit      int64 `json:"limit" form:"limit"`
	Pagination bool  `json:"is_paginated" form:"is_paginated"`
}

func (p *Pagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 25
	}
	if p.Page == 0 {
		p.Page = 1
	}
	p.Pagination = true
}
