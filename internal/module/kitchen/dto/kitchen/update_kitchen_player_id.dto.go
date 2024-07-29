package kitchen

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UpdateKitchenPlayerIdDto struct {
	PlayerId string `json:"player_id" validate:"required"`
	dto.KitchenHeaders
}

func (payload *UpdateKitchenPlayerIdDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
