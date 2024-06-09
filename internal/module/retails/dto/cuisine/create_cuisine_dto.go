package cuisine

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3,max=30"`
	En string `json:"en" validate:"required,min=3,max=30"`
}

type CreateCuisineDto struct {
	Name     Name   `json:"name" validate:"required"`
	Logo     string `json:"logo"`
	IsHidden bool   `json:"is_hidden"`
}

func (input *CreateCuisineDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
