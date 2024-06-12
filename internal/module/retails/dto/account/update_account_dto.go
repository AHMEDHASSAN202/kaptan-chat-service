package account

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateAccountDto struct {
	Name            Name     `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email,Account_Email_is_unique_rules_validation"`
	Password        string   `json:"password"`
	Country         Country  `json:"country" validate:"required"`
	AllowedBrandIds []string `json:"allowed_brand_ids" validate:"required"`
}

func (payload *UpdateAccountDto) Validate(c echo.Context, validate *validator.Validate, validateAccountEmailIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Account_Email_is_unique_rules_validation,
		RegisterValidationFunc: validateAccountEmailIsUnique,
	})
}
