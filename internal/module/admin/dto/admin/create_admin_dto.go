package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" bson:"ar" validate:"required"`
	En string `json:"en" bson:"en" validate:"required"`
}

type Account struct {
	Id   string `json:"id" validate:"required"`
	Name Name   `json:"name" validate:"required"`
}
type Kitchen struct {
	Id            string   `json:"id" validate:"required"`
	Name          Name     `json:"name" validate:"required"`
	LocationIds   []string `json:"location_ids" validate:"required"`
	AccountIds    []string `json:"account_ids" validate:"required"`
	AllowedStatus []string `json:"allowed_status"`
}

type CreateAdminDTO struct {
	ID              primitive.ObjectID `json:"-"`
	Name            string             `json:"name" validate:"required,min=3"`
	Email           string             `json:"email" validate:"required,Email_is_unique_rules_validation"`
	Status          string             `json:"status" validate:"oneof=active inactive"`
	Password        string             `json:"password" validate:"Password_required_if_id_is_zero,omitempty,min=8"`
	ConfirmPassword string             `json:"password_confirmation" validate:"required_with=Password,eqfield=Password"`
	Type            string             `json:"type" validate:"required,oneof=admin portal kitchen"`
	RoleId          string             `json:"role_id" validate:"required,mongodb,RoleExistsValidation"`
	CountryIds      []string           `json:"country_ids" validate:"required,country_ids"`
	AdminDetails    dto.AdminDetails   `json:"-"`
	Account         *Account           `json:"account"`
	Kitchen         *Kitchen           `json:"kitchen"`
	dto.AdminHeaders
}

func (input *CreateAdminDTO) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool, passwordRequiredIfIdIsZero func(fl validator.FieldLevel) bool, roleExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	validate.RegisterValidation("country_ids", utils.ValidateCountryIds)
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateEmailIsUnique,
	},
		validators.CustomErrorTags{
			ValidationTag:          localization.Password_required_if_id_is_zero,
			RegisterValidationFunc: passwordRequiredIfIdIsZero,
		},
		validators.CustomErrorTags{
			ValidationTag:          localization.RoleExistsValidation,
			RegisterValidationFunc: roleExists,
		})
}
