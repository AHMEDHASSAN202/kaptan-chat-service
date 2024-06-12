package item

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ChangeItemStatusDto struct {
	Id           string             `json:"_"`
	Status       string             `json:"status" validate:"required,oneof=active inactive"`
	AdminDetails []dto.AdminDetails `json:"-"`
}

func (input *ChangeItemStatusDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input)
}
