package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListOrderDtoForDashboard struct {
	dto.Pagination
	Status      string `query:"status" validate:"omitempty,oneof=initiated pending accepted ready_for_pickup pickedup no_show cancelled rejected timeout"`
	Query       string `json:"query" form:"query" query:"query"`
	IsFavourite bool   `query:"is_favourite"`
	From        string `query:"from" validate:"omitempty,DateTimeFormat"`
	To          string `query:"to" validate:"omitempty,DateTimeFormat"`
	AccountId   string `query:"account_id"`
	BrandId     string `query:"brand_id"`
	LocationId  string `query:"location_id"`
	dto.AdminHeaders
}

func (input *ListOrderDtoForDashboard) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}

type ListOrderDtoForMobile struct {
	dto.Pagination
	UserId string `header:"causer-id" validate:"required"`
}

func (input *ListOrderDtoForMobile) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}

type FindOrderMobileDto struct {
	UserId  string `header:"causer-id" validate:"required"`
	OrderId string `json:"order_id"`
}

func (payload *FindOrderMobileDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
