package responses

import (
	"github.com/kamva/mgm/v3"
)

type MobileListCuisine struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name   `json:"name" validate:"required,dive"`
	Logo             string `json:"logo" bson:"logo"`
	IsHidden         bool   `json:"is_hidden" bson:"is_hidden"`
}
