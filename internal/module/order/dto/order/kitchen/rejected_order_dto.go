package kitchen

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type RejectedOrderDto struct {
	OrderId          string `param:"id" validate:"required,mongodb"`
	RejectedReasonId string `json:"rejected_reason_id" validate:"required"`
	Note             string `json:"note"`
	dto.KitchenHeaders
}

func (d *RejectedOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
