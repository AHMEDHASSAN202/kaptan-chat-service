package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type ToggleOrderFavDto struct {
	UserId  string `header:"causer-id" validate:"required"`
	OrderId string `json:"order_id"`
}

func (payload *ToggleOrderFavDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
