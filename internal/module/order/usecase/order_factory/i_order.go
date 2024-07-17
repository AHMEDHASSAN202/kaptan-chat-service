package order_factory

import (
	"context"
	"samm/internal/module/order/domain"
	kitchen3 "samm/internal/module/order/responses/kitchen"
	"samm/internal/module/order/responses/user"
	"samm/pkg/validators"
)

type IOrder interface {
	Make() IOrder
	Create(ctx context.Context, dto interface{}) (*user.FindOrderResponse, validators.ErrorResponse)
	Find(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse)
	List(ctx context.Context, dto interface{}) ([]domain.Order, validators.ErrorResponse)
	SendNotifications() validators.ErrorResponse
	ToPending(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse)
	ToAcceptKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse)
	ToRejectedKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse)
	ToReadyForPickupKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse)
	ToPickedUpKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse)
	ToNoShowKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse)
	ToArrived(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse)
	ToCancel(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse)
}
