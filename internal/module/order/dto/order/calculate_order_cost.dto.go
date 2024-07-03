package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type MenuItem struct {
	Id          string     `json:"id" validate:"required"`
	Qty         int64      `json:"qty" validate:"required,min=1"`
	ModifierIds []MenuItem `json:"modifier_ids"`
}

type CalculateOrderCostDto struct {
	LocationId         string     `json:"location_id" validate:"required,mongodb"`
	AccountId          string     `json:"account_id" validate:"required,mongodb"`
	UserId             string     `json:"-"`
	CollectionMethodId string     `json:"collection_method_id" validate:"required"`
	MenuItems          []MenuItem `json:"menu_items" validate:"required,dive"`
}

func (d *CalculateOrderCostDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
