package modifier_group

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DeleteModifierGroupDto struct {
	Id string `param:"id" validate:"required"`
	dto.PortalHeaders
}

func (input *DeleteModifierGroupDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
