package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateKitchenProfileDTO struct {
	dto.PortalHeaders
	ID              primitive.ObjectID `json:"-"`
	Name            string             `json:"name" validate:"required,min=3"`
	Email           string             `json:"email" validate:"required,Email_is_unique_rules_validation"`
	Password        string             `json:"password" validate:"Password_required_if_id_is_zero,omitempty,min=8"`
	ConfirmPassword string             `json:"password_confirmation" validate:"required_with=Password,eqfield=Password"`
	AdminDetails    dto.AdminDetails   `json:"-"`
}

func (input *UpdateKitchenProfileDTO) Validate(c echo.Context, validate *validator.Validate, validateEmailIsUnique func(fl validator.FieldLevel) bool, passwordRequiredIfIdIsZero func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Email_is_unique_rules_validation,
		RegisterValidationFunc: validateEmailIsUnique,
	},
		validators.CustomErrorTags{
			ValidationTag:          localization.Password_required_if_id_is_zero,
			RegisterValidationFunc: passwordRequiredIfIdIsZero,
		})
}
