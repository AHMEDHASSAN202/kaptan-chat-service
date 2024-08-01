package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type MissedItems struct {
	Id            string        `json:"id" validate:"required"`
	Qty           int64         `json:"qty" validate:"required,min=1"`
	MissingAddons []MissedItems `json:"missing_addons"`
}

type ReportMissingItemDto struct {
	OrderId      string        `param:"id" validate:"required"`
	MissingItems []MissedItems `json:"missing_items" validate:"required"`
}

func (payload *ReportMissingItemDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
