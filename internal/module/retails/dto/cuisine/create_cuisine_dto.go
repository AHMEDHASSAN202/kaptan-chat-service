package cuisine

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3,max=30,Cuisine_name_is_unique_rules_validation"`
	En string `json:"en" validate:"required,min=3,max=30,Cuisine_name_is_unique_rules_validation"`
}

type CreateCuisineDto struct {
	Name     Name   `json:"name" validate:"required"`
	Logo     string `json:"logo"`
	IsHidden bool   `json:"is_hidden"`
	dto.AdminHeaders
}

func (input *CreateCuisineDto) Validate(c echo.Context, validate *validator.Validate, validateCuisineNameExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_name_is_unique_rules_validation,
		RegisterValidationFunc: validateCuisineNameExists,
	})
}
