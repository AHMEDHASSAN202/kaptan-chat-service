package kitchen

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type AcceptOrderDto struct {
	PreparationTime int `json:"preparation_time" validate:"required"`
	dto.KitchenHeaders
}

func (d *AcceptOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
