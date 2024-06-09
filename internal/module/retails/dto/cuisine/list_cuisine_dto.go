package cuisine

import "samm/pkg/utils/dto"

type ListCuisinesDto struct {
	dto.Pagination
	Query string `json:"query,omitempty" form:"query,omitempty"`
}
