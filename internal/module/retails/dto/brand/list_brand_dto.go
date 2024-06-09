package brand

import "samm/pkg/utils/dto"

type ListBrandDto struct {
	dto.Pagination
	Query string `json:"query,omitempty" form:"query,omitempty"`
}
