package structs

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Type string `json:"type" bson:"type"`
	Num  string `json:"num" bson:"num"`
}

type OrderPriceSummary struct {
	Fees                     float64 `json:"fees" bson:"fees"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount" bson:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount" bson:"total_price_after_discount"`
}

type CollectionMethod struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Type   string             `json:"type" bson:"type"`
	Fields map[string]any     `json:"fields" bson:"fields"`
	Values map[string]any     `json:"values" bson:"values"`
}

type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Name             string             `json:"name" bson:"name"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	Country          string             `json:"country" bson:"country"`
	CollectionMethod CollectionMethod   `json:"collection_method" bson:"collection_method"`
}

type MobileFindOrder struct {
	mgm.DefaultModel `bson:",inline"`
	Status           string            `json:"status" bson:"status"`
	Location         Location          `json:"location" bson:"location"`
	SerialNum        string            `json:"serial_num" bson:"serial_num"`
	Items            []Item            `json:"items" bson:"items"`
	IsFavourite      bool              `json:"is_favourite" bson:"is_favourite"`
	Notes            string            `json:"notes" bson:"notes"`
	User             User              `json:"user" bson:"user"`
	Payment          Payment           `json:"payment" bson:"payment"`
	PriceSummary     OrderPriceSummary `json:"price_summary" bson:"price_summary"`
}
