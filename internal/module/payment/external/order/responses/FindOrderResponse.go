package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type Country struct {
	Id          string           `json:"id" bson:"_id"`
	Name        LocalizationText `json:"name" bson:"name"`
	Timezone    string           `json:"timezone" bson:"timezone"`
	Currency    string           `json:"currency" bson:"currency"`
	PhonePrefix string           `json:"phone_prefix" bson:"phone_prefix"`
}
type Location struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            LocalizationText   `json:"name" bson:"name"`
	Street          LocalizationText   `json:"street" bson:"street"`
	CoverImage      string             `json:"cover_image" bson:"cover_image"`
	PreparationTime int                `json:"preparation_time" bson:"preparation_time"`
	Logo            string             `json:"logo" bson:"logo"`
	Phone           string             `json:"phone" bson:"phone"`
	Percent         float64            `json:"percent" bson:"percent"`
	Country         Country            `json:"country" bson:"country"`
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
	CollectionMethod *CollectionMethod  `json:"collection_method" bson:"collection_method"`
}

type OrderPriceSummary struct {
	Fees                     float64 `json:"fees" bson:"fees"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount" bson:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount" bson:"total_price_after_discount"`
}
type OrderResponse struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	SerialNum        string             `json:"serial_num" bson:"serial_num"`
	User             User               `json:"user" bson:"user"`
	Location         Location           `json:"location" bson:"location"`
	PreparationTime  int                `json:"preparation_time" bson:"preparation_time"`
	PriceSummary     OrderPriceSummary  `json:"price_summary" bson:"price_summary"`
	Status           string             `json:"status" bson:"status"`
	IsFavourite      bool               `json:"is_favourite" bson:"is_favourite"`
	AcceptedAt       *time.Time         `json:"accepted_at" bson:"accepted_at"`
	PaidAt           *time.Time         `json:"paid_at" bson:"paid_at"`
	ArrivedAt        *time.Time         `json:"arrived_at" bson:"arrived_at"`
	PickedUpAt       *time.Time         `json:"pickedup_at" bson:"pickedup_at"`
	ReadyForPickUpAt *time.Time         `json:"ready_for_pickup_at" bson:"ready_for_pickup_at"`
	CancelledAt      *time.Time         `json:"cancelled_at" bson:"cancelled_at"`
	RejectedAt       *time.Time         `json:"rejected_at" bson:"rejected_at"`
	NoShowAt         *time.Time         `json:"no_show_at" bson:"no_show_at"`
	Notes            string             `json:"notes" bson:"notes"`
}
