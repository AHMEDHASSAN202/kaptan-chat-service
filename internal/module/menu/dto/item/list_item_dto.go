package item

import "samm/pkg/utils/dto"

type ListItemsDto struct {
	dto.Pagination
	Query     string `json:"query" form:"query"`
	AccountId string `json:"account_id" form:"account_id"`
}
