package item

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ItemAvailability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3,Item_name_is_unique_rules_validation"`
	En string `json:"en" validate:"required,min=3,Item_name_is_unique_rules_validation"`
}

type LocalizationTextDesc struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type ItemResponse struct {
	mgm.DefaultModel `bson:",inline"`
	AccountId        primitive.ObjectID       `json:"account_id" bson:"account_id"`
	Name             LocalizationText         `json:"name" bson:"name"`
	Desc             LocalizationText         `json:"desc" bson:"desc"`
	Type             string                   `json:"type" bson:"type"`
	Min              int                      `json:"min" bson:"min"`
	Max              int                      `json:"max" bson:"max"`
	Calories         int                      `json:"calories" bson:"calories"`
	Price            float64                  `json:"price" bson:"price"`
	ModifierGroupIds []map[string]interface{} `json:"modifier_groups_ids" bson:"modifier_groups_ids"`
	Availabilities   []ItemAvailability       `json:"availabilities" bson:"availabilities"`
	Tags             []string                 `json:"tags" bson:"tags"`
	Image            string                   `json:"image" bson:"image"`
	AdminDetails     []map[string]interface{} `json:"admin_details" bson:"admin_details"`
	Status           string                   `json:"status" bson:"status"`
	DeletedAt        *time.Time               `json:"deleted_at" bson:"deleted_at"`
}
