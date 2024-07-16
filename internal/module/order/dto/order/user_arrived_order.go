package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ArrivedOrderDto struct {
	OrderId            string
	UserId             string `header:"causer-id" validate:"required"`
	CollectionMethodId string `json:"collection_method_id"`

	dto.MobileHeaders
}

func (d *ArrivedOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
