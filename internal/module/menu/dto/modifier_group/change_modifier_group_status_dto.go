package modifier_group

import "samm/pkg/utils/dto"

type ChangeModifierGroupStatusDto struct {
	Id           string             `json:"_"`
	Status       string             `json:"status"`
	AdminDetails []dto.AdminDetails `json:"-"`
}
