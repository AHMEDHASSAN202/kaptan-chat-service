package location

import "samm/pkg/utils/dto"

type ListLocationMobileDto struct {
	dto.Pagination
	dto.MobileHeaders
	CountryId string `header:"Country-Id"`
	Query     string `query:"query"`
	BrandId   string `query:"brand_id"`
}
