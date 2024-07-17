package kitchen

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type PickedUpOrderDto struct {
	OrderId string `param:"id" validate:"required,mongodb"`
	dto.KitchenHeaders
}

func (d *PickedUpOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
