package account

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type BankAccount struct {
	AccountNumber string `json:"account_number" validate:"required"`
	BankName      string `json:"bank_name" validate:"required"`
	CompanyName   string `json:"company_name" validate:"required"`
}

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
	Name            Name        `json:"name" validate:"required"`
	Email           string      `json:"email" validate:"required,email,Email_is_unique_rules_validation"`
	Password        string      `json:"password" validate:"required,omitempty,min=8"`
	ConfirmPassword string      `json:"password_confirmation" validate:"required,eqfield=Password"`
	Country         Country     `json:"country" validate:"required"`
	AllowedBrandIds []string    `json:"allowed_brand_ids" validate:"required,min=1"`
	Percent         float64     `json:"percent" validate:"required"`
	BankAccount     BankAccount `json:"bank_account" validate:"required"`
}

func (payload *StoreAccountDto) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateEmailIsUnique,
	})
}
