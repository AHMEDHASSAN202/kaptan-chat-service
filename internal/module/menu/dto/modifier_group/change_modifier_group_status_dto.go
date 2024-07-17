package modifier_group

import (
	"context"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
)

type ChangeModifierGroupStatusDto struct {
	Id           string             `param:"id" validate:"required"`
	Status       string             `json:"status" validate:"required,oneof=active inactive"`
	AdminDetails []dto.AdminDetails `json:"-"`
	dto.PortalHeaders
}

func (input *ChangeModifierGroupStatusDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
