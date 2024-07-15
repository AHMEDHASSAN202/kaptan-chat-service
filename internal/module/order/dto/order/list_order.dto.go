package order

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListOrderDto struct {
	dto.Pagination
	Status     string `query:"status" validate:"omitempty,oneof=initiated pending accepted ready_for_pickup pickedup no_show cancelled rejected timeout"`
	From       string `query:"from" validate:"omitempty,DateTimeFormat"`
	To         string `query:"to" validate:"omitempty,DateTimeFormat"`
	AccountId  string `query:"account_id"`
	BrandId    string `query:"brand_id"`
	SerialNum  string `query:"serial_num"`
	CountryId  string `query:"country_id"`
	LocationId string `query:"location_id"`
	UserId     string `query:"user_id"`
}

func (input *ListOrderDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
