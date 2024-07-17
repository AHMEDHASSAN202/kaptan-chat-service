package item

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DeleteItemDto struct {
	Id string `json:"_"`
	dto.PortalHeaders
}

func (input *DeleteItemDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
