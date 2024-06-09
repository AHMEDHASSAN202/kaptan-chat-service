package account

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type UpdateAccountDto struct {
	Name     Name   `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
}

func (payload *UpdateAccountDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
