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

type CreateAdminDTO struct {
	ID              primitive.ObjectID `json:"-"`
	Name            string             `json:"name" validate:"required,min=3"`
	Email           string             `json:"email" validate:"required,Email_is_unique_rules_validation"`
	Status          string             `json:"status" validate:"oneof=active inactive"`
	Password        string             `json:"password" validate:"Password_required_if_id_is_zero,omitempty,min=8"`
	ConfirmPassword string             `json:"password_confirmation" validate:"required_with=Password,eqfield=Password"`
	Type            string             `json:"type" validate:"required,oneof=admin portal"`
	Role            string             `json:"role" validate:"required"`
	Permissions     []string           `json:"permissions" validate:"required"`
	CountryIds      []string           `json:"country_ids" validate:"required,country_ids"`
	AccountId       string             `json:"account_id" validate:"required_if=Type portal"`
	AdminDetails    dto.AdminDetails   `json:"-"`
}

func (input *CreateAdminDTO) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool, passwordRequiredIfIdIsZero func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	validate.RegisterValidation("country_ids", utils.ValidateCountryIds)
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateEmailIsUnique,
	},
		validators.CustomErrorTags{
			ValidationTag:          localization.Password_required_if_id_is_zero,
			RegisterValidationFunc: passwordRequiredIfIdIsZero,
		})
}
