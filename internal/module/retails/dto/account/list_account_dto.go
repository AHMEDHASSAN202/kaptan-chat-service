package account

import "samm/pkg/utils/dto"

type ListAccountDto struct {
	dto.Pagination
	Query string `query:"query"`
}
