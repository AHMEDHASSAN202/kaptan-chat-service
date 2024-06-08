package item

import "samm/pkg/utils/dto"

type ChangeItemStatusDto struct {
	Id           string             `json:"_"`
	Status       string             `json:"status"`
	AdminDetails []dto.AdminDetails `json:"-"`
}
