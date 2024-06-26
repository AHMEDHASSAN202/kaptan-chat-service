package location

import "samm/pkg/utils/dto"

type ListLocationMobileDto struct {
	dto.Pagination
	dto.MobileHeaders
	CountryId  string   `header:"Country-Id"`
	Query      string   `query:"query"`
	BrandId    string   `query:"brand_id"`
	Distance   float64  `query:"distance"`
	CuisineIds []string `query:"cuisine_ids"`
}

func (p *ListLocationMobileDto) SetDefault() {
	p.Pagination.SetDefault()
	if p.Distance == 0 {
		p.Distance = 20
	}
}
