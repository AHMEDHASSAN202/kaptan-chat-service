package Notification

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type RedirectData struct {
	LocationId primitive.ObjectID `json:"location_id" bson:"location_id"`
}

type NotificationMobile struct {
	mgm.DefaultModel `bson:",inline"`
	Title            Name          `json:"title" bson:"title"`
	Image            string        `json:"image" bson:"image"`
	Text             Name          `json:"text" bson:"text"`
	RedirectType     string        `json:"redirect_type" bson:"redirect_type"`
	RedirectData     *RedirectData `json:"redirect_data" bson:"redirect_data"`
}
