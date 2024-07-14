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

type Order struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name       `json:"name" bson:"name"`
	Email            string     `json:"email" bson:"email"`
	Password         string     `json:"-" bson:"password"`
	DeletedAt        *time.Time `json:"-" bson:"deleted_at"`
}

type OrderV2 struct {
	mgm.DefaultModel `bson:",inline"`
	SerialNum        string `json:"serial_num"`
	User             struct {
		Id               string `json:"id"`
		Name             string `json:"name"`
		Phone            string `json:"phone"`
		Country          string `json:"country"`
		CollectionMethod struct {
		} `json:"collection_method"`
	} `json:"user"` // need to update
	Items []struct {
		ItemDetails struct {
			Id     string `json:"_id"`
			Name   string `json:"name"`
			Qty    string `json:"qty"`
			Price  string `json:"price"`
			Addons struct {
				Id    string `json:"_id"`
				Name  string `json:"name"`
				Qty   string `json:"qty"`
				Price string `json:"price"`
			} `json:"addons"`
			Category struct {
			} `json:"category"`
		} `json:"item_details"`
	} `json:"items"` // need to update
	LocationDetails struct {
	} `json:"location_details"` // need to update
	PreparationTime string     `json:"preparation_time"`
	TotalPrice      float64    `json:"total_price"`
	OrderStatus     string     `json:"order_status"`
	IsFavourite     bool       `json:"is_favourite"`
	CreatedAt       *time.Time `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	PaidAt          *time.Time `json:"paid_at"`
	PickedupAt      *time.Time `json:"pickedup_at"`
	ReadyForPickAt  *time.Time `json:"ready_for_pick_at"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	RejectedAt      *time.Time `json:"rejected_at"`
	Cancelled       struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"cancelled"`
	Rejected struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"rejected"`
	StatusLog []struct {
		Id     string `json:"id"`
		Type   string `json:"type"`
		Status struct {
			New string `json:"new"`
			Old string `json:"old"`
		} `json:"status"`
		CreatedAt *time.Time `json:"created_at"`
	} `json:"status_log"`
	Notes          string `json:"notes"`
	PaymentDetails struct {
		Type string `json:"type"`
		Num  string `json:"num"`
	} `json:"payment_details"`
}
type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, payload *order.StoreOrderDto) (err validators.ErrorResponse)
	CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse)
	UpdateOrder(ctx context.Context, id string, payload *order.UpdateOrderDto) (err validators.ErrorResponse)
	FindOrder(ctx context.Context, Id string) (order Order, err validators.ErrorResponse)
	DeleteOrder(ctx context.Context, Id string) (err validators.ErrorResponse)
	ListOrderForDashboard(ctx context.Context, payload *order.ListOrderDto) (*responses.ListResponse, validators.ErrorResponse)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (err error)
	UpdateOrder(ctx context.Context, order *Order) (err error)
	FindOrder(ctx context.Context, Id primitive.ObjectID) (order *Order, err error)
	DeleteOrder(ctx context.Context, Id primitive.ObjectID) (err error)
	ListOrderForDashboard(ctx *context.Context, dto *order.ListOrderDto) (ordersRes *[]OrderV2, paginationMeta *PaginationData, err error)
}
