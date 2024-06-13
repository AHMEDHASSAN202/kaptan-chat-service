package account

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}
type Country struct {
	Id   string `json:"_id" validate:"required"`
	Name struct {
		Ar string `json:"ar" validate:"required"`
		En string `json:"en" validate:"required"`
	} `json:"name" validate:"required"`
	Timezone    string `json:"timezone" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
}

type StoreAccountDto struct {
	Name            Name     `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email,Account_Email_is_unique_rules_validation"`
	Password        string   `json:"password" validate:"required"`
	ConfirmPassword string   `json:"password_confirmation" validate:"required,eqfield=Password"`
	Country         Country  `json:"country" validate:"required"`
	AllowedBrandIds []string `json:"allowed_brand_ids" validate:"required,min=1"`
}

func (payload *StoreAccountDto) Validate(c echo.Context, validate *validator.Validate, validateAccountEmailIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Account_Email_is_unique_rules_validation,
		RegisterValidationFunc: validateAccountEmailIsUnique,
	})
}
