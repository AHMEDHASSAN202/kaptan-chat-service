package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
)

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type Doc struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  Name               `json:"name" bson:"name"`
	Image string             `json:"image" bson:"image"`
}

type Account struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name Name               `json:"name" bson:"name"`
}

type CreateApprovalDto struct {
	dto.AdminDetails
	CountryId  string
	EntityId   primitive.ObjectID
	EntityType string
	New        map[string]interface{}
	Old        map[string]interface{}
	Doc        Doc
	Account    Account
}
