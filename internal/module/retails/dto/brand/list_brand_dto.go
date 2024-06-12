package brand

import "samm/pkg/utils/dto"

type ListBrandDto struct {
	dto.Pagination
	Query string   `json:"query" query:"query"`
	Ids   []string `json:"ids" query:"ids"`
}
