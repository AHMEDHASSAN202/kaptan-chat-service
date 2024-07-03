package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type LoginAsPortalDto struct {
	Id           string             `param:"id" validate:"required,mongodb"`
	AdminDetails []dto.AdminDetails `json:"-"`
	dto.AdminHeaders
}

func (input *LoginAsPortalDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
