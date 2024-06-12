package menu_group_item

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"time"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type ItemAvailability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}
type MenuGroupItemSyncItemModel struct {
	UpdatedAt        time.Time            `json:"updated_at" bson:"updated_at"`
	ItemId           primitive.ObjectID   `json:"item_id" bson:"item_id"`
	AccountId        primitive.ObjectID   `json:"account_id" bson:"account_id"`
	Name             LocalizationText     `json:"name" bson:"name"`
	Desc             LocalizationText     `json:"desc" bson:"desc"`
	Calories         int                  `json:"calories" bson:"calories"`
	Price            float64              `json:"price" bson:"price"`
	ModifierGroupIds []primitive.ObjectID `json:"modifier_group_ids" bson:"modifier_group_ids"`
	Availabilities   []ItemAvailability   `json:"availabilities" bson:"availabilities"`
	Tags             []string             `json:"tags" bson:"tags"`
	Image            string               `json:"image" bson:"image"`
	AdminDetails     []dto.AdminDetails   `json:"admin_details" bson:"admin_details,omitempty"`
	Status           string               `json:"status" bson:"status"`
}
