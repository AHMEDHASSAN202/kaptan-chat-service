package location

import "samm/pkg/utils/dto"

type ListLocationDto struct {
	dto.Pagination
	Query     string `query:"query"`
	AccountId string `query:"account_id"`
	BrandId   string `query:"brand_id"`
}
