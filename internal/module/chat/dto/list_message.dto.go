package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type ListChatDtoForDashboard struct {
	dto.Pagination
	Status      string `query:"status" validate:"omitempty,oneof=initiated pending accepted ready_for_pickup pickedup no_show cancelled rejected timeout"`
	Query       string `json:"query" form:"query" query:"query"`
	IsFavourite bool   `query:"is_favourite"`
	From        string `query:"from" validate:"omitempty,DateTimeFormat"`
	To          string `query:"to" validate:"omitempty,DateTimeFormat"`
	AccountId   string `query:"account_id"`
	BrandId     string `query:"brand_id"`
	LocationId  string `query:"location_id"`
}

func (input *ListChatDtoForDashboard) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}

type ListChatDtoForMobile struct {
	dto.Pagination
	UserId string `header:"causer-id" validate:"required"`
}

type FindOrderMobileDto struct {
	UserId  string `header:"causer-id" validate:"required"`
	OrderId string `json:"order_id"`
}

func (payload *FindOrderMobileDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
