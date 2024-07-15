package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/responses"
	"samm/pkg/validators"
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

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
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
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	ItemId       primitive.ObjectID `json:"item_id" bson:"item_id"`
	Name         LocalizationText   `json:"name" bson:"name"`
	Desc         LocalizationText   `json:"desc" bson:"desc"`
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

type City struct {
	Id   primitive.ObjectID `json:"id" bson:"id"`
	Name LocalizationText   `json:"name" bson:"name"`
}

type Brand struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name LocalizationText   `json:"name" bson:"name"`
	Logo string             `json:"logo" bson:"logo"`
}

type PercentsDate struct {
	From    time.Time `json:"from" bson:"from"`
	To      time.Time `json:"to" bson:"to"`
	Percent float64   ` json:"percent" bson:"percent"`
}

type Country struct {
	Id   string `json:"id" bson:"_id"`
	Name struct {
		Ar string `json:"ar" bson:"ar"`
		En string `json:"en" bson:"en"`
	} `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}

type Location struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            LocalizationText   `json:"name" bson:"name"`
	City            City               `json:"city" bson:"city"`
	Street          LocalizationText   `json:"street" bson:"street"`
	CoverImage      string             `json:"cover_image" bson:"cover_image"`
	PreparationTime int                `json:"preparation_time" bson:"preparation_time"`
	Logo            string             `json:"logo" bson:"logo"`
	Phone           string             `json:"phone" bson:"phone"`
	Brand           Brand              `json:"brand_details" bson:"brand_details"`
	Percent         float64            `json:"percent" bson:"percent"`
	PercentsDate    []PercentsDate     `json:"percents_date" bson:"percents_date"`
	Country         Country            `json:"country" bson:"country"`
	AccountId       primitive.ObjectID `json:"account_id" bson:"account_id"`
}

type Rejected struct {
	Id       string `json:"id" bson:"id"`
	Note     string `json:"note" bson:"note"`
	UserType string `json:"user_type" bson:"user_type"`
}

type StatusLog struct {
	CauserId   string `json:"causer_id" bson:"causer_id"`
	CauserType string `json:"causer_type" bson:"causer_type"`
	Status     struct {
		New string `json:"new" bson:"new"`
		Old string `json:"old" bson:"old"`
	} `json:"status" bson:"status"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
}

type Payment struct {
	Type string `json:"type" bson:"type"`
	Num  string `json:"num" bson:"num"`
}

type Order struct {
	mgm.DefaultModel `bson:",inline"`
	SerialNum        string            `json:"serial_num" bson:"serial_num"`
	User             User              `json:"user" bson:"user"`
	Items            []Item            `json:"items" bson:"items"`
	Location         Location          `json:"location" bson:"location"`
	PreparationTime  int               `json:"preparation_time" bson:"preparation_time"`
	PriceSummary     OrderPriceSummary `json:"price_summary" bson:"price_summary"`
	Status           string            `json:"status" bson:"status"`
	IsFavourite      bool              `json:"is_favourite" bson:"is_favourite"`
	AcceptedAt       *time.Time        `json:"accepted_at" bson:"accepted_at"`
	PaidAt           *time.Time        `json:"paid_at" bson:"paid_at"`
	PickedUpAt       *time.Time        `json:"pickedup_at" bson:"pickedup_at"`
	ReadyForPickUpAt *time.Time        `json:"ready_for_pickup_at" bson:"ready_for_pickup_at"`
	CancelledAt      *time.Time        `json:"cancelled_at" bson:"cancelled_at"`
	RejectedAt       *time.Time        `json:"rejected_at" bson:"rejected_at"`
	NoShowAt         *time.Time        `json:"no_show_at" bson:"no_show_at"`
	Cancelled        *Rejected         `json:"cancelled,omitempty" bson:"cancelled,omitempty"`
	Rejected         *Rejected         `json:"rejected,omitempty" bson:"rejected,omitempty"`
	StatusLogs       []StatusLog       `json:"status_logs" bson:"status_logs"`
	Notes            string            `json:"notes" bson:"notes"`
	Payment          Payment           `json:"payment" bson:"payment"`
}

type OrderUseCase interface {
	CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse)
	ListOrderForDashboard(ctx context.Context, payload *order.ListOrderDto) (*responses.ListResponse, validators.ErrorResponse)
	StoreOrder(ctx context.Context, payload *order.CreateOrderDto) (interface{}, validators.ErrorResponse)

	UserRejectionReasons(ctx context.Context, status string, id string) ([]UserRejectionReason, validators.ErrorResponse)

	UserCancelOrder(ctx context.Context, payload *order.CancelOrderDto) (Order, validators.ErrorResponse)
}

type OrderRepository interface {
	StoreOrder(ctx *context.Context, order *Order) (err error)
	ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDto) (ordersRes *[]Order, paginationMeta *PaginationData, err error)
}
