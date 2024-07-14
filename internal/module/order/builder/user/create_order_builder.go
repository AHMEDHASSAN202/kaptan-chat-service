package user

import (
	"context"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	responses2 "samm/internal/module/order/external/menu/responses"
	"samm/internal/module/order/external/retails/responses"
	"samm/pkg/validators"
)

func CreateOrderBuilder(ctx context.Context, dto *order.CreateOrderDto, location responses.LocationDetails, items []responses2.MenuDetailsResponse) (*domain.Order, validators.ErrorResponse) {
	orderModel := domain.Order{}

	return &orderModel, validators.ErrorResponse{}
}
