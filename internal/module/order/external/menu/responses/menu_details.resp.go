package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type MenuDetailsResponse struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	ItemId           primitive.ObjectID   `json:"item_id" bson:"item_id"`
	Name             LocalizationText     `json:"name" bson:"name"`
	Desc             LocalizationText     `json:"desc" bson:"desc"`
	Calories         int                  `json:"calories" bson:"calories"`
	Price            float64              `json:"price" bson:"price"`
	ModifierGroupIds []primitive.ObjectID `json:"modifier_group_ids" bson:"modifier_group_ids"`
	ModifierGroups   []ModifierGroup      `json:"modifier_groups" bson:"modifier_groups"`
	Addons           []MobileGetItemAddon `json:"addons" bson:"addons"`
	Category         GetItemCategory      `json:"category" bson:"category"`
	Tags             []string             `json:"tags" bson:"tags"`
	Image            string               `json:"image" bson:"image"`
}

type GetItemCategory struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name LocalizationText   `json:"name" bson:"name"`
	Icon string             `json:"icon" bson:"icon"`
}

type ModifierGroup struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id"`
	Name       LocalizationText     `json:"name" bson:"name"`
	Type       string               `json:"type" bson:"type"`
	Min        int                  `json:"min" bson:"min"`
	Max        int                  `json:"max" bson:"max"`
	ProductIds []primitive.ObjectID `json:"product_ids" bson:"product_ids"`
}

type MobileGetItemAddon struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     LocalizationText   `json:"name" bson:"name"`
	Type     string             `json:"type" bson:"type"`
	Min      int                `json:"min" bson:"min"`
	Max      int                `json:"max" bson:"max"`
	Calories int                `json:"calories" bson:"calories"`
	Price    float64            `json:"price" bson:"price"`
	Image    string             `json:"image" bson:"image"`
}
