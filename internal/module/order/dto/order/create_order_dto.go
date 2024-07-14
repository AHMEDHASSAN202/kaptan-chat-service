package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type CreateOrderDto struct {
	LocationId         string     `json:"location_id" validate:"required,mongodb"`
	UserId             string     `header:"causer-id" validate:"required"`
	CollectionMethodId string     `json:"collection_method_id"`
	MenuItems          []MenuItem `json:"menu_items" validate:"required,dive"`
}

func (d *CreateOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
