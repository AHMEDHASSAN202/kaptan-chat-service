package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type PortalAuthDTO struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

func (input *PortalAuthDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
