package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/responses"
	"time"
)

type CollectionMethod struct {
	ID     primitive.ObjectID `json:"id"`
	Type   string             `json:"type"`
	Fields map[string]any     `json:"fields"`
	Values map[string]any     `json:"values"`
}

type User struct {
	ID               primitive.ObjectID `json:"id"`
	Name             string             `json:"name"`
	PhoneNumber      string             `json:"phone_number"`
	Country          string             `json:"country"`
	CollectionMethod CollectionMethod   `json:"collection_method"`
}

type ItemPriceSummary struct {
	Qty                      int     `json:"qty"`
	UnitPrice                float64 `json:"unit_price"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount"`
}

type OrderPriceSummary struct {
	Fees                     float64 `json:"fees"`
	TotalPriceBeforeDiscount float64 `json:"total_price_before_discount"`
	TotalPriceAfterDiscount  float64 `json:"total_price_after_discount"`
}

type Item struct {
	ID           primitive.ObjectID         `json:"id"`
	Name         responses.LocalizationText `json:"name"`
	Desc         responses.LocalizationText `json:"desc"`
	Type         string                     `json:"type"`
	Min          int                        `json:"min"`
	Max          int                        `json:"max"`
	SKU          string                     `json:"sku"`
	Calories     int                        `json:"calories"`
	Price        float64                    `json:"price"`
	Image        string                     `json:"image"`
	Qty          int                        `json:"qty"`
	PriceSummary ItemPriceSummary           `json:"price_summary"`
	Addons       []Item                     `json:"addons"`
}

type City struct {
	Id   primitive.ObjectID         `json:"id"`
	Name responses.LocalizationText `json:"name"`
}

type Brand struct {
	ID   primitive.ObjectID         `json:"id"`
	Name responses.LocalizationText `json:"name"`
	Logo string                     `json:"logo"`
}

type Country struct {
	ID   string `json:"id"`
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name"`
	Timezone    string `json:"timezone"`
	Currency    string `json:"currency"`
	PhonePrefix string `json:"phone_prefix"`
}
type Account struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name struct {
		Ar string `json:"ar"`
		En string `json:"en"`
	} `json:"name" bson:"name"`
	AllowedBrandIds []primitive.ObjectID `json:"allowed_brand_ids" bson:"allowed_brand_ids"`
}
type Location struct {
	ID              primitive.ObjectID         `json:"id"`
	Name            responses.LocalizationText `json:"name"`
	City            City                       `json:"city"`
	Street          responses.LocalizationText `json:"street"`
	CoverImage      string                     `json:"cover_image"`
	PreparationTime int                        `json:"preparation_time"`
	Logo            string                     `json:"logo"`
	Phone           string                     `json:"phone"`
	Brand           Brand                      `json:"brand_details"`
	Country         Country                    `json:"country"`
	Account         Account                    `json:"account" bson:"account"`
}

type Rejected struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Payment struct {
	Type string `json:"type"`
	Num  string `json:"num"`
}

type FindOrderResponse struct {
	ID               primitive.ObjectID `json:"id"`
	SerialNum        string             `json:"serial_num"`
	User             User               `json:"user"`
	Items            []Item             `json:"items"`
	Location         Location           `json:"location"`
	PreparationTime  int                `json:"preparation_time"`
	PriceSummary     OrderPriceSummary  `json:"price_summary"`
	Status           string             `json:"status"`
	IsFavourite      bool               `json:"is_favourite"`
	AcceptedAt       *time.Time         `json:"accepted_at"`
	PaidAt           *time.Time         `json:"paid_at"`
	PickedUpAt       *time.Time         `json:"pickedup_at"`
	ReadyForPickUpAt *time.Time         `json:"ready_for_pickup_at"`
	CancelledAt      *time.Time         `json:"cancelled_at"`
	RejectedAt       *time.Time         `json:"rejected_at"`
	NoShowAt         *time.Time         `json:"no_show_at"`
	Cancelled        *Rejected          `json:"cancelled,omitempty"`
	Rejected         *Rejected          `json:"rejected,omitempty"`
	Notes            string             `json:"notes"`
	Payment          Payment            `json:"payment"`
}
