package domain

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils"
	"time"
)

type Status struct {
	SnoozeTo string `json:"snooze_to"`
	Status   string `json:"status" bson:"status"`
	Meta     Meta   `json:"meta" bson:"meta"`
}
type Meta struct {
	NameEN string `json:"name_en" bson:"name_en"`
	NameAr string `json:"name_ar" bson:"name_ar"`
	Color  string `json:"color" bson:"color"`
}

type LocationMobile struct {
	mgm.DefaultModel           `bson:",inline"`
	Name                       Name               `json:"name" bson:"name"`
	City                       City               `json:"city" bson:"city"`
	Street                     Name               `json:"street" bson:"street"`
	CoverImage                 string             `json:"cover_image" bson:"cover_image"`
	Logo                       string             `json:"logo" bson:"logo"`
	SnoozeTo                   *time.Time         `json:"snooze_to" bson:"snooze_to"`
	IsOpen                     bool               `json:"is_open" bson:"is_open"`
	WorkingHour                []WorkingHour      `json:"working_hour" bson:"working_hour"`
	Phone                      string             `json:"phone" bson:"phone"`
	Coordinate                 Coordinate         `json:"coordinate" bson:"coordinate"`
	BrandDetails               BrandDetails       `json:"brand_details" bson:"brand_details"`
	PreparationTime            int                `json:"preparation_time" bson:"preparation_time"`
	Distance                   float64            `json:"distance" bson:"distance"`
	Country                    Country            `json:"country" bson:"country"`
	AutoAccept                 bool               `json:"auto_accept" bson:"auto_accept"`
	Status                     Status             `json:"status" bson:"-"`
	AccountId                  primitive.ObjectID `json:"account_id" bson:"account_id"`
	AllowedCollectionMethodIds []string           `json:"-" bson:"allowed_collection_method_ids"`
	AllowedCollectionMethods   []interface{}      `json:"allowed_collection_methods,omitempty"`
}

func (payload *LocationMobile) SetOpenStatus() {
	open := payload.IsOpen
	now := time.Now().UTC()

	if open {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "open",
			Meta: Meta{
				NameEN: "Open",
				NameAr: "مفتوح",
				Color:  "",
			},
		}
	} else {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "closed",
			Meta: Meta{
				NameEN: "Closed",
				NameAr: "مغلق",
				Color:  "",
			},
		}
	}

	if payload.SnoozeTo != nil && now.Before(*payload.SnoozeTo) {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "busy",
			Meta: Meta{
				NameEN: "Busy",
				NameAr: "مشغول",
				Color:  "",
			},
		}
	}

}
func (payload *LocationMobile) SetDistance(lat float64, lng float64) {
	if lat != 0 && lng != 0 {
		payload.Distance = utils.Distance(payload.Coordinate.Coordinates[1], payload.Coordinate.Coordinates[0], lat, lng)
	}
}
