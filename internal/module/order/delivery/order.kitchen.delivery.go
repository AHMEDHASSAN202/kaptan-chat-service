package delivery

import (
	"github.com/labstack/echo/v4"
	"samm/internal/module/order/dto/order/kitchen"
	"samm/pkg/validators"
)

func (a *OrderHandler) KitchenToAccept(c echo.Context) error {
	ctx := c.Request().Context()

	var acceptDto kitchen.AcceptOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &acceptDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&acceptDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := acceptDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	orderResponse, errResp := a.orderUsecase.AcceptOrder(ctx, &acceptDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
