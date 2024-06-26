package card

import (
	"samm/pkg/utils/dto"
)

type ListCardDto struct {
	dto.Pagination
	UserId string `json:"user_id" `
}
