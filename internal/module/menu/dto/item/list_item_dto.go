package item

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListItemsDto struct {
	dto.Pagination
	Query     string `json:"query" query:"query"`
	AccountId string `json:"account_id" query:"account_id" header:"account_id" validate:"required"`
}

func (input *ListItemsDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
