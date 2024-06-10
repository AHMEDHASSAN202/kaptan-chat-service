package modifier_group

import "samm/pkg/utils/dto"

type ListModifierGroupsDto struct {
	dto.Pagination
	Query string `json:"query" form:"query" query:"query"`
}
