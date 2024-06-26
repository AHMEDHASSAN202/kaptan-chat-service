package modifier_group

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"time"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type ModifierGroupResp struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name         LocalizationText     `json:"name" bson:"name"`
	Type         string               `json:"type" bson:"type"`
	Min          int                  `json:"min" bson:"min"`
	Max          int                  `json:"max" bson:"max"`
	Products     []map[string]any     `json:"products" bson:"products"`
	ProductIds   []primitive.ObjectID `json:"product_ids" bson:"product_ids"`
	Status       string               `json:"status" bson:"status"`
	AdminDetails []dto.AdminDetails   `json:"admin_details" bson:"admin_details"`
	AccountId    primitive.ObjectID   `json:"account_id" bson:"account_id"`
	DeletedAt    *time.Time           `json:"deleted_at" bson:"deleted_at"`
}
