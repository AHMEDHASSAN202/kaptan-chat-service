package user

import (
	"samm/pkg/utils/dto"
)

type ListUserDto struct {
	dto.Pagination
	Query  string `query:"query"`
	Status bool   `query:"status,omitempty"`
	dob    string `query:"dob"`
}
