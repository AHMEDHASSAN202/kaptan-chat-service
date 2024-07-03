package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/responses"
	"samm/pkg/utils"
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

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type OrderUseCase interface {
	StoreOrder(ctx context.Context, payload *order.StoreOrderDto) (err validators.ErrorResponse)
	CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse)
	UpdateOrder(ctx context.Context, id string, payload *order.UpdateOrderDto) (err validators.ErrorResponse)
	FindOrder(ctx context.Context, Id string) (order Order, err validators.ErrorResponse)
	DeleteOrder(ctx context.Context, Id string) (err validators.ErrorResponse)
	ListOrder(ctx context.Context, payload *order.ListOrderDto) (orders []Order, paginationResult utils.PaginationResult, err validators.ErrorResponse)
}

type OrderRepository interface {
	StoreOrder(ctx context.Context, order *Order) (err error)
	UpdateOrder(ctx context.Context, order *Order) (err error)
	FindOrder(ctx context.Context, Id primitive.ObjectID) (order *Order, err error)
	DeleteOrder(ctx context.Context, Id primitive.ObjectID) (err error)
	ListOrder(ctx context.Context, payload *order.ListOrderDto) (locations []Order, paginationResult utils.PaginationResult, err error)
}
