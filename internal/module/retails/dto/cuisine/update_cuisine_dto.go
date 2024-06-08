package cuisine

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type UpdateCuisineDto struct {
	Id string `json:"_"`
	CreateCuisineDto
}

func (input *UpdateCuisineDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
