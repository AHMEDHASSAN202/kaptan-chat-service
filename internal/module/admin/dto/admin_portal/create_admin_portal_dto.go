package admin_portal

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type Account struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name Name               `json:"name" bson:"name"`
}

type CreateAdminPortalDTO struct {
	ID              primitive.ObjectID `json:"-"`
	Name            string             `json:"name" validate:"required,min=3"`
	Email           string             `json:"email" validate:"required,Email_is_unique_rules_validation"`
	Status          string             `json:"status" validate:"oneof=active inactive"`
	Password        string             `json:"password" validate:"Password_required_if_id_is_zero,omitempty,min=8"`
	ConfirmPassword string             `json:"password_confirmation" validate:"required_with=Password,eqfield=Password"`
	RoleId          string             `json:"role_id" validate:"required,mongodb,RoleExistsValidation"`
	Account         *Account           `json:"account" validate:"AccountRequiredValidation"`
	AdminDetails    dto.AdminDetails   `json:"-"`
	dto.PortalHeaders
}

func (input *CreateAdminPortalDTO) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool, passwordRequiredIfIdIsZero func(fl validator.FieldLevel) bool, roleExists func(fl validator.FieldLevel) bool, validation_account func(fl validator.FieldLevel) bool) validators.ErrorResponse {
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
		},
		validators.CustomErrorTags{
			ValidationTag:          localization.AccountRequiredValidation,
			RegisterValidationFunc: validation_account,
		})
}
