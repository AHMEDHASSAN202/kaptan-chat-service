package brand

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3,max=30"`
	En string `json:"en" validate:"required,min=3,max=30"`
}

type CreateBrandDto struct {
	Name       Name     `json:"name" validate:"required"`
	Logo       string   `json:"logo"`
	IsActive   bool     `json:"is_active"`
	CuisineIds []string `json:"cuisine_ids" validate:"Invalid_mongo_ids_validation_rule,Cuisine_id_is_exists_rules_validation"`
	AccountId  string   `json:"account_id" validate:"omitempty,Account_id_is_not_exists"`
	dto.AdminHeaders
}

func (input *CreateBrandDto) Validate(ctx context.Context, validate *validator.Validate, validateCuisineExists func(fl validator.FieldLevel) bool, validateAccountExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_id_is_exists_rules_validation,
		RegisterValidationFunc: validateCuisineExists,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Account_id_is_not_exists,
		RegisterValidationFunc: validateAccountExists,
	})
}

func existsIModifierGroup(db interface{}, value []string) bool {
	return true
}
