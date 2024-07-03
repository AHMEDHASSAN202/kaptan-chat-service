package role

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type DeleteRoleDTO struct {
	ID string `param:"id" validate:"required,mongodb,PreventDeleteRolesIdsValidation,RoleHasAdminsValidation"`
}

func (input *DeleteRoleDTO) Validate(c echo.Context, validate *validator.Validate, roleHasAdminsValidation func(fl validator.FieldLevel) bool, preventDeleteRolesIdsValidation func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStructAndReturnOneError(c.Request().Context(), validate, input,
		validators.CustomErrorTags{
			ValidationTag:          localization.RoleHasAdminsValidation,
			RegisterValidationFunc: roleHasAdminsValidation,
		},
		validators.CustomErrorTags{
			ValidationTag:          localization.PreventDeleteRolesIdsValidation,
			RegisterValidationFunc: preventDeleteRolesIdsValidation,
		},
	)
}
