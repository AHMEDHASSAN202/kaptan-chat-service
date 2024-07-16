package structs

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type Brand struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name LocalizationText   `json:"name" bson:"name"`
	Logo string             `json:"logo" bson:"logo"`
}
type City struct {
	Id   primitive.ObjectID `json:"id" bson:"id"`
	Name LocalizationText   `json:"name" bson:"name"`
}

type Location struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       LocalizationText   `json:"name" bson:"name"`
	City       City               `json:"city" bson:"city"`
	Street     LocalizationText   `json:"street" bson:"street"`
	CoverImage string             `json:"cover_image" bson:"cover_image"`
	Logo       string             `json:"logo" bson:"logo"`
	Brand      Brand              `json:"brand_details" bson:"brand_details"`
	AccountId  primitive.ObjectID `json:"account_id" bson:"account_id"`
}

type ItemPriceSummary struct {
	Qty                      int     `json:"qty" bson:"qty"`
	UnitPrice                float64 `json:"unit_price" bson:"unit_price"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount" bson:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount" bson:"total_price_after_discount"`
}

type Item struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	ItemId       primitive.ObjectID `json:"item_id" bson:"item_id"`
	Name         LocalizationText   `json:"name" bson:"name"`
	Type         string             `json:"type" bson:"type,omitempty"`
	Min          int                `json:"min" bson:"min,omitempty"`
	Max          int                `json:"max" bson:"max,omitempty"`
	SKU          string             `json:"sku" bson:"sku"`
	Calories     int                `json:"calories" bson:"calories"`
	Price        float64            `json:"price" bson:"price"`
	Image        string             `json:"image" bson:"image"`
	Qty          int                `json:"qty" bson:"qty"`
	PriceSummary ItemPriceSummary   `json:"price_summary" bson:"price_summary"`
	Addons       []Item             `json:"addons" bson:"addons,omitempty"`
}

type MobileListOrders struct {
	mgm.DefaultModel `bson:",inline"`
	Status           string   `json:"status" bson:"status"`
	Location         Location `json:"location" bson:"location"`
	SerialNum        string   `json:"serial_num" bson:"serial_num"`
	Items            []Item   `json:"items" bson:"items"`
	IsFavourite      bool     `json:"is_favourite" bson:"is_favourite"`
}
