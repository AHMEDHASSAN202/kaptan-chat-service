package kitchen

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type FindOrderMobileDto struct {
	KitchenId string `header:"causer-id" validate:"required"`
	OrderId   string `param:"id"`
}

func (payload *FindOrderMobileDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
