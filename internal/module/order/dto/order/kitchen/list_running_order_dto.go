package kitchen

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListRunningOrderDto struct {
	Status             []string `json:"status" query:"status[]" validate:"required"`
	NumberOfHoursLimit int      `json:"-" query:"-"`
	dto.KitchenHeaders
	dto.Pagination
}

func (d *ListRunningOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
