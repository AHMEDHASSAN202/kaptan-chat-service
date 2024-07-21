package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DashboardCancelOrderDto struct {
	OrderId string
	Note    string `json:"note"`
	dto.AdminHeaders
}

func (d *DashboardCancelOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
