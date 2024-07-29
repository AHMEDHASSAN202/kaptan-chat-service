package brand

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateBrandDto struct {
	Id string `json:"_"`
	CreateBrandDto
}

func (input *UpdateBrandDto) Validate(ctx context.Context, validate *validator.Validate, validateCuisineExists func(fl validator.FieldLevel) bool, validateAccountExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_id_is_exists_rules_validation,
		RegisterValidationFunc: validateCuisineExists,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Account_id_is_not_exists,
		RegisterValidationFunc: validateAccountExists,
	})
}
