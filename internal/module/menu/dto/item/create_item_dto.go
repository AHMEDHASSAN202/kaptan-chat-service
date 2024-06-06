package item

import "samm/pkg/utils/dto"

type ItemAvailability struct {
	Day  string `json:"day"`
	From string `json:"from"`
	To   string `json:"to"`
}
type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type LocalizationTextDesc struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}
type CreateItemDto struct {
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
