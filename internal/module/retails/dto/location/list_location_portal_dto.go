package location

import "samm/pkg/utils/dto"

type ListLocationPortalDto struct {
	dto.Pagination
	Query   string `query:"query"`
	BrandId string `query:"brand_id"`
	dto.PortalHeaders
}
