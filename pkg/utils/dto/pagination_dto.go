package dto

type Pagination struct {
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
	Sort   string `json:"sort,omitempty;query:sort"`
	SortBy string `json:"sort_by,omitempty;query:sort"`
}

func (p *Pagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 25
	}
	if p.Page == 0 {
		p.Page = 1
	}
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
