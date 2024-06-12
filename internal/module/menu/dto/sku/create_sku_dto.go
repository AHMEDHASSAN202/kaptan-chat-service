package sku

import (
	"context"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"

	"github.com/go-playground/validator/v10"
)

type CreateSKUDto struct {
	Name string `json:"name" validate:"required,SKU_name_is_unique_rules_validation"`
}

func (input *CreateSKUDto) Validate(ctx context.Context, validate *validator.Validate, validateNameIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.SKU_name_is_unique_rules_validation,
		RegisterValidationFunc: validateNameIsUnique,
	})
}
