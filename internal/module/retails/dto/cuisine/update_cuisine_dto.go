package cuisine

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateCuisineDto struct {
	Id string `json:"_"`
	CreateCuisineDto
}

func (input *UpdateCuisineDto) Validate(c echo.Context, validate *validator.Validate, validateCuisineNameExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_name_is_unique_rules_validation,
		RegisterValidationFunc: validateCuisineNameExists,
	})
}
