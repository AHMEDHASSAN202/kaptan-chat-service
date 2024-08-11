package payment

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UpdateSession struct {
	SessionId string `json:"session_id" validate:"required"`
	dto.MobileHeaders
}

func (payload *UpdateSession) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
