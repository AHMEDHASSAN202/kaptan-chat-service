package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

type LocalizationText struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type AccountByIdResp struct {
	ID   primitive.ObjectID `json:"id"`
	Name LocalizationText   `json:"name"`
}
