package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type CancelOrderDto struct {
	OrderId        string
	UserId         string `header:"causer-id" validate:"required"`
	CancelReasonId string `json:"cancel_reason_id" validate:"required"`
	dto.MobileHeaders
}

func (d *CancelOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
