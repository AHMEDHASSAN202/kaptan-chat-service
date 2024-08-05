package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/dto/order/kitchen"
	"samm/internal/module/order/repository/structs"
	"samm/internal/module/order/responses"
	"samm/internal/module/order/responses/user"
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
	CollectionMethod *CollectionMethod  `json:"collection_method" bson:"collection_method"`
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
type MissedItem struct {
	Id  string `json:"id,omitempty" bson:"id,omitempty"`
	Qty int64  `json:"qty,omitempty" bson:"qty,omitempty"`
}
type Item struct {
	ID               primitive.ObjectID  `json:"id" bson:"_id"`
	ItemId           primitive.ObjectID  `json:"item_id" bson:"item_id"`
	MobileId         string              `json:"mobile_id" bson:"mobile_id"`
	Name             LocalizationText    `json:"name" bson:"name"`
	Desc             LocalizationText    `json:"desc" bson:"desc"`
	Type             string              `json:"type" bson:"type,omitempty"`
	Min              int                 `json:"min" bson:"min,omitempty"`
	Max              int                 `json:"max" bson:"max,omitempty"`
	SKU              string              `json:"sku" bson:"sku"`
	Calories         int                 `json:"calories" bson:"calories"`
	Price            float64             `json:"price" bson:"price"`
	Image            string              `json:"image" bson:"image"`
	Qty              int                 `json:"qty" bson:"qty"`
	PriceSummary     ItemPriceSummary    `json:"price_summary" bson:"price_summary"`
	ModifierGroupId  *primitive.ObjectID `json:"modifier_group_id" bson:"modifier_group_id,omitempty"`
	Addons           []Item              `json:"addons" bson:"addons,omitempty"`
	MissedItemReport *MissedItem         `json:"missed_item_report,omitempty" bson:"missed_item_report,omitempty"`
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
type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}
type Country struct {
	Id          string `json:"id" bson:"_id"`
	Name        Name   `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}
type Account struct {
	Id              primitive.ObjectID   `json:"id" bson:"_id"`
	Name            Name                 `json:"name" bson:"name"`
	AllowedBrandIds []primitive.ObjectID `json:"allowed_brand_ids" bson:"allowed_brand_ids"`
}
type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
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
	Account         Account            `json:"account" bson:"account"`
	Coordinate      Coordinate         `json:"coordinate" bson:"coordinate"`
}

type Rejected struct {
	Id       string `json:"id" bson:"id"`
	Note     string `json:"note" bson:"note"`
	Name     *Name  `json:"name" bson:"name"`
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
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	PaymentType string             `json:"payment_type" bson:"payment_type"`
	CardType    string             `json:"card_type" bson:"card_type"`
	CardNumber  string             `json:"card_number" bson:"card_number"`
}

type MetaData struct {
	HasMissingItems  bool                 `json:"has_missing_items" bson:"has_missing_items"`
	TargetKitchenIds []primitive.ObjectID `bson:"target_kitchen_ids" json:"target_kitchen_ids"`
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
	ArrivedAt        *time.Time        `json:"arrived_at" bson:"arrived_at"`
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
	MetaData         MetaData          `json:"meta_data" bson:"meta_data"`
}

type OrderUseCase interface {
	StoreOrder(ctx context.Context, payload *order.CreateOrderDto) (interface{}, validators.ErrorResponse)
	ReportMissedItem(ctx context.Context, payload *order.ReportMissingItemDto) (interface{}, validators.ErrorResponse)
	KitchenAcceptOrder(ctx context.Context, payload *kitchen.AcceptOrderDto) (interface{}, validators.ErrorResponse)
	KitchenRejectedOrder(ctx context.Context, payload *kitchen.RejectedOrderDto) (interface{}, validators.ErrorResponse)
	KitchenPickedUpOrder(ctx context.Context, payload *kitchen.PickedUpOrderDto) (interface{}, validators.ErrorResponse)
	KitchenNoShowOrder(ctx context.Context, payload *kitchen.NoShowOrderDto) (interface{}, validators.ErrorResponse)
	KitchenReadyForPickupOrder(ctx context.Context, payload *kitchen.ReadyForPickupOrderDto) (interface{}, validators.ErrorResponse)
	KitchenRejectionReasons(ctx context.Context, status string, id string) ([]KitchenRejectionReason, validators.ErrorResponse)
	CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse)
	ListOrderForDashboard(ctx context.Context, payload *order.ListOrderDtoForDashboard) (*responses.ListResponse, validators.ErrorResponse)
	ListInprogressOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse)
	ListCompletedOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse)
	ListLastOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse)
	FindOrderForDashboard(ctx *context.Context, id string) (*Order, validators.ErrorResponse)
	FindOrderForMobile(ctx *context.Context, payload *order.FindOrderMobileDto) (*user.FindOrderResponse, validators.ErrorResponse)
	ToggleOrderFavourite(ctx *context.Context, payload order.ToggleOrderFavDto) (err validators.ErrorResponse)
	DashboardCancelOrder(ctx context.Context, payload *order.DashboardCancelOrderDto) (*Order, validators.ErrorResponse)
	DashboardPickedOrder(ctx context.Context, payload *order.DashboardPickedUpOrderDto) (*Order, validators.ErrorResponse)

	UserRejectionReasons(ctx context.Context, status string, id string) ([]UserRejectionReason, validators.ErrorResponse)

	UpdateRealTimeDb(ctx context.Context, order *Order) validators.ErrorResponse

	UserCancelOrder(ctx context.Context, payload *order.CancelOrderDto) (*user.FindOrderResponse, validators.ErrorResponse)
	UserArrivedOrder(ctx context.Context, payload *order.ArrivedOrderDto) (*user.FindOrderResponse, validators.ErrorResponse)
	SetOrderPaid(ctx context.Context, payload *order.OrderPaidDto) validators.ErrorResponse

	// cron jobs
	CronJobTimedOutOrders(ctx context.Context) validators.ErrorResponse
	CronJobPickedOrders(ctx context.Context) validators.ErrorResponse
	CronJobCancelOrders(ctx context.Context) validators.ErrorResponse
}

type OrderRepository interface {
	StoreOrder(ctx context.Context, order *Order) (*Order, error)
	UpdateOrder(ctx context.Context, order *Order) (err error)
	FindOrder(ctx *context.Context, Id primitive.ObjectID) (*Order, error)
	ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDtoForDashboard) (ordersRes *[]Order, paginationMeta *PaginationData, err error)
	ListInprogressOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error)
	ListCompletedOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error)
	ListLastOrdersForMobile(ctx *context.Context, dto *order.ListOrderDtoForMobile) (ordersRes *[]structs.MobileListOrders, paginationMeta *PaginationData, err error)
	UserHasOrders(ctx context.Context, userId primitive.ObjectID, orderStatus []string, gt int64) (bool, error)
	FindOrderByUser(ctx *context.Context, id string, userId string) (order *Order, err error)
	UpdateOrderStatus(ctx *context.Context, orderDomain *Order, previousStatus []string, statusLog *StatusLog, updateSet interface{}) (order *Order, err error)
	UpdateUserAllOrdersFavorite(ctx context.Context, userId string) (err error)
	GetAllOrdersForCronJobs(ctx *context.Context, filters bson.M) (ordersRes *[]Order, paginationMeta *PaginationData, err error)
}
