package item

import "samm/pkg/utils/dto"

type UpdateItemDto struct {
	Id                string               `json:"_"`
	AccountId         string               `json:"account_id"`
	Name              LocalizationText     `json:"name"`
	Desc              LocalizationTextDesc `json:"desc"`
	Type              string               `json:"type"`
	Min               int                  `json:"min"`
	Max               int                  `json:"max"`
	Calories          int                  `json:"calories"`
	Price             float64              `json:"price"`
	ModifierGroupsIds []string             `json:"modifier_groups_ids"`
	Availabilities    []ItemAvailability   `json:"availabilities"`
	Tags              []string             `json:"tags"`
	Image             string               `json:"image"`
	Status            string               `json:"status"`
	AdminDetails      []dto.AdminDetails   `json:"-"`
}
