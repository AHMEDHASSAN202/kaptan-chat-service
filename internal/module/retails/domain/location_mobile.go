package domain

import (
	"github.com/kamva/mgm/v3"
)

type LocationMobile struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name   `json:"name" bson:"name"`
	City             City   `json:"city" bson:"city"`
	Street           Name   `json:"street" bson:"street"`
	CoverImage       string `json:"cover_image" bson:"cover_image"`
	Logo             string `json:"logo" bson:"logo"`
	// Open Status
	Phone           string       `json:"phone" bson:"phone"`
	Coordinate      Coordinate   `json:"coordinate" bson:"coordinate"`
	BrandDetails    BrandDetails `json:"brand_details" bson:"brand_details"`
	PreparationTime int          `json:"preparation_time" bson:"preparation_time"`
	Country         Country      `json:"country" bson:"country"`
}
