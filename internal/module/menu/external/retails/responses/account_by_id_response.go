package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	Id   string `json:"_id" bson:"_id"`
	Name struct {
		Ar string `json:"ar" bson:"ar"`
		En string `json:"en" bson:"en"`
	} `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}

type AccountByIdResp struct {
	ID      primitive.ObjectID `json:"id"`
	Name    LocalizationText   `json:"name"`
	Country Country            `json:"country"`
}
