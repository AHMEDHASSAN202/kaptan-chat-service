package brand

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3,max=30"`
	En string `json:"en" validate:"required,min=3,max=30"`
}

type CreateBrandDto struct {
	Name        Name     `json:"name" validate:"required"`
	Logo        string   `json:"logo"`
	IsActive    bool     `json:"is_active"`
	SnoozedTill string   `json:"snoozed_till"`
	CuisineIds  []string `json:"cuisine_ids" validate:"Invalid_mongo_ids_validation_rule,Cuisine_id_is_exists_rules_validation"`
}

func (input *CreateBrandDto) Validate(ctx context.Context, validate *validator.Validate, validateCuisineExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_id_is_exists_rules_validation,
		RegisterValidationFunc: validateCuisineExists,
	})
}

func existsIModifierGroup(db interface{}, value []string) bool {
	return true
}
