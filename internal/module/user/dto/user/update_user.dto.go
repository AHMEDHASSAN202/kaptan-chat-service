package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateUserProfileDto struct {
	Name string `json:"name" validate:"required"`
	//CountryCode string `json:"country_code" validate:"required,len=4,numeric"`
	//PhoneNumber string `json:"phone_number" validate:"required,PhoneNumber_rule_validation"`
	Email  string `json:"email" validate:"required,email,Email_is_unique_rules_validation"`
	Gender string `json:"gender" validate:"required,oneof=male female other"`
	Dob    string `json:"dob" validate:"required,datetime=2006-01-02"`
	//ImageURL string `json:"image_url" validate:"omitempty,url"`
	dto.MobileHeaders
}

func (payload *UpdateUserProfileDto) Validate(ctx context.Context, validate *validator.Validate, validateUserEmailIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateUserEmailIsUnique,
	})
}
