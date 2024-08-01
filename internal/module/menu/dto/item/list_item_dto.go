package item

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListItemsDto struct {
	dto.Pagination
	Query       string `json:"query" query:"query"`
	Type        string `json:"type" query:"type" validate:"omitempty,oneof=product modifier"`
	AccountId   string `json:"account_id" query:"account_id" header:"account_id" validate:"required"`
	HasOriginal *bool  `json:"has_original" query:"has_original"`
}

func (input *ListItemsDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
