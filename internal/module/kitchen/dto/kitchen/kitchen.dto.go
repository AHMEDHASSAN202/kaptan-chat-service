package kitchen

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}
type Country struct {
	Id   string `json:"_id"  validate:"required"`
	Name struct {
		Ar string `json:"ar"  validate:"required"`
		En string `json:"en"  validate:"required"`
	} `json:"name"  validate:"required"`
	Timezone    string `json:"timezone"  validate:"required"`
	Currency    string `json:"currency"  validate:"required"`
	PhonePrefix string `json:"phone_prefix"  validate:"required"`
}

type StoreKitchenDto struct {
	ID              primitive.ObjectID `json:"-"`
	Name            Name               `json:"name" validate:"required"`
	Email           string             `json:"email" validate:"required,email,Email_is_unique_rules_validation"`
	Password        string             `json:"password" validate:"required,omitempty,min=8"`
	ConfirmPassword string             `json:"password_confirmation" validate:"required,eqfield=Password"`
	Country         Country            `json:"country" validate:"required"`

	AccountIds    []string `json:"account_ids" validate:"required,Validate_Account_Location_validation"`
	LocationIds   []string `json:"location_ids" validate:"required,Validate_Account_Location_validation"`
	AllowedStatus []string `json:"allowed_status" validate:"required"`

	dto.AdminHeaders
}

func (payload *StoreKitchenDto) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool, validateAccountAndLocation func(fl validator.FieldLevel) bool) validators.ErrorResponse {

	return validators.ValidateStruct(c.Request().Context(), validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateEmailIsUnique,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Validate_Account_Location_validation,
		RegisterValidationFunc: validateAccountAndLocation,
	})
}

type ListKitchenDto struct {
	dto.Pagination
	Query      string `query:"query"`
	AccountId  string `query:"account_id"`
	LocationId string `query:"location_id"`
}

type UpdateKitchenDto struct {
	Name    Name    `json:"name" validate:"required"`
	Country Country `json:"country" validate:"required"`

	AccountIds  []string `json:"account_ids" validate:"required,Validate_Account_Location_validation"`
	LocationIds []string `json:"location_ids" validate:"required,Validate_Account_Location_validation"`
	dto.AdminHeaders
}

func (payload *UpdateKitchenDto) Validate(c echo.Context, validate *validator.Validate, validateAccountAndLocation func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload, validators.CustomErrorTags{
		ValidationTag:          localization.Validate_Account_Location_validation,
		RegisterValidationFunc: validateAccountAndLocation,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Validate_Account_Location_validation,
		RegisterValidationFunc: validateAccountAndLocation,
	})
}

type DeleteKitchenDto struct {
	Id string
	dto.AdminHeaders
}
