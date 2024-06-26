package user

import (
	"samm/pkg/utils/dto"
)

type ListUserDto struct {
	dto.Pagination
	Query string `query:"query"`
}
