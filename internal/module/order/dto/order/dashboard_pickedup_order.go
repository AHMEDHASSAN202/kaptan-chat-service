package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DashboardPickedUpOrderDto struct {
	OrderId string
	dto.AdminHeaders
}

func (d *DashboardPickedUpOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
