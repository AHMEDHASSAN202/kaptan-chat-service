package kitchen

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/responses"
	"time"
)

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

type ItemPriceSummary struct {
	Qty                      int     `json:"qty" bson:"qty"`
	UnitPrice                float64 `json:"unit_price" bson:"unit_price"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount" bson:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount" bson:"total_price_after_discount"`
}

type OrderPriceSummary struct {
	Fees                     float64 `json:"fees" bson:"fees"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount" bson:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount" bson:"total_price_after_discount"`
}

type Item struct {
	ID           primitive.ObjectID         `json:"id" bson:"_id"`
	Name         responses.LocalizationText `json:"name" bson:"name"`
	MobileId     string                     `json:"mobile_id,omitempty" bson:"mobile_id"`
	Desc         responses.LocalizationText `json:"desc" bson:"desc"`
	Type         string                     `json:"type" bson:"type"`
	Min          int                        `json:"min" bson:"min"`
	Max          int                        `json:"max" bson:"max"`
	SKU          string                     `json:"sku" bson:"sku"`
	Calories     int                        `json:"calories" bson:"calories"`
	Price        float64                    `json:"price" bson:"price"`
	Image        string                     `json:"image" bson:"image"`
	Qty          int                        `json:"qty" bson:"qty"`
	PriceSummary ItemPriceSummary           `json:"price_summary" bson:"price_summary"`
	Addons       []Item                     `json:"addons" bson:"addons"`
}

type City struct {
	Id   primitive.ObjectID         `json:"id" bson:"id"`
	Name responses.LocalizationText `json:"name" bson:"name"`
}

type Brand struct {
	ID   primitive.ObjectID         `json:"id" bson:"_id"`
	Name responses.LocalizationText `json:"name" bson:"name"`
	Logo string                     `json:"logo" bson:"logo"`
}

type Country struct {
	ID   string `json:"id" bson:"_id"`
	Name struct {
		Ar string `json:"ar" bson:"ar"`
		En string `json:"en" bson:"en"`
	} `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}

type Location struct {
	ID              primitive.ObjectID         `json:"id" bson:"_id"`
	Name            responses.LocalizationText `json:"name" bson:"name"`
	City            City                       `json:"city" bson:"city"`
	Street          responses.LocalizationText `json:"street" bson:"street"`
	CoverImage      string                     `json:"cover_image" bson:"cover_image"`
	PreparationTime int                        `json:"preparation_time" bson:"preparation_time"`
	Logo            string                     `json:"logo" bson:"logo"`
	Phone           string                     `json:"phone" bson:"phone"`
	Brand           Brand                      `json:"brand_details" bson:"brand_details"`
	Country         Country                    `json:"country" bson:"country"`
	Account         Account                    `json:"account" bson:"account"`
}

type Account struct {
	Id              primitive.ObjectID   `json:"id" bson:"_id"`
	Name            Name                 `json:"name" bson:"name"`
	AllowedBrandIds []primitive.ObjectID `json:"allowed_brand_ids" bson:"allowed_brand_ids"`
}
type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type Rejected struct {
	Id       string                     `json:"id" bson:"_id"`
	Note     string                     `json:"note" bson:"note"`
	Name     responses.LocalizationText `json:"name" bson:"name"`
	UserType string                     `json:"user_type" bson:"user_type"`
}

type Payment struct {
	Type string `json:"type" bson:"type"`
	Num  string `json:"num" bson:"num"`
}

type FindOrderResponse struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	SerialNum        string             `json:"serial_num" bson:"serial_num"`
	User             User               `json:"user" bson:"user"`
	Items            []Item             `json:"items" bson:"items"`
	Location         Location           `json:"location" bson:"location"`
	PreparationTime  int                `json:"preparation_time" bson:"preparation_time"`
	PriceSummary     OrderPriceSummary  `json:"price_summary" bson:"price_summary"`
	Status           string             `json:"status" bson:"status"`
	IsFavourite      bool               `json:"is_favourite" bson:"is_favourite""`
	AcceptedAt       *time.Time         `json:"accepted_at" bson:"accepted_at"`
	PaidAt           *time.Time         `json:"paid_at" bson:"paid_at"`
	PickedUpAt       *time.Time         `json:"pickedup_at" bson:"pickedup_at"`
	ReadyForPickUpAt *time.Time         `json:"ready_for_pickup_at" bson:"ready_for_pickup_at"`
	CancelledAt      *time.Time         `json:"cancelled_at" bson:"cancelled_at"`
	RejectedAt       *time.Time         `json:"rejected_at" bson:"rejected_at"`
	NoShowAt         *time.Time         `json:"no_show_at" bson:"no_show_at"`
	Cancelled        *Rejected          `json:"cancelled,omitempty" bson:"cancelled"`
	Rejected         *Rejected          `json:"rejected,omitempty" bson:"rejected"`
	Notes            string             `json:"notes" bson:"notes"`
	Payment          Payment            `json:"payment" bson:"payment"`
}
