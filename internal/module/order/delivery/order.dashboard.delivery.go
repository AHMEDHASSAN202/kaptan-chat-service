package delivery

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/internal/module/order/dto/order"
	"samm/pkg/validators"
)

func (a *OrderHandler) FindOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	order, errResp := a.orderUsecase.FindOrderForDashboard(&ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, order)
}

func (a *OrderHandler) ListOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDtoForDashboard

	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	payload.Pagination.SetDefault()

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	orders, errResp := a.orderUsecase.ListOrderForDashboard(ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, orders)

}

func (a *OrderHandler) CancelOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	var orderDto order.DashboardCancelOrderDto

	orderId := c.Param("id")

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	orderDto.OrderId = orderId

	validationErr := orderDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	orderResponse, errResp := a.orderUsecase.DashboardCancelOrder(ctx, &orderDto)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		errResp.StatusCode = http.StatusBadRequest
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
func (a *OrderHandler) PickedUpOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	var orderDto order.DashboardPickedUpOrderDto

	orderId := c.Param("id")

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	orderDto.OrderId = orderId

	validationErr := orderDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	orderResponse, errResp := a.orderUsecase.DashboardPickedOrder(ctx, &orderDto)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		errResp.StatusCode = http.StatusBadRequest
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
