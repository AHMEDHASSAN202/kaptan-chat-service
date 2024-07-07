package modifier_group

import "samm/pkg/utils/dto"

type ListModifierGroupsDto struct {
	dto.Pagination
	Query     string `json:"query" form:"query" query:"query"`
	AccountId string `json:"account_id" query:"account_id" header:"account_id" validate:"required"`
}
