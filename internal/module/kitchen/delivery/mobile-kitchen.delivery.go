package delivery

import (
	"github.com/labstack/echo/v4"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/pkg/validators"
)

func (a *KitchenHandler) UpdateKitchenPlayerId(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.UpdateKitchenPlayerIdDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.kitchenUsecase.UpdateKitchenPlayerId(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
