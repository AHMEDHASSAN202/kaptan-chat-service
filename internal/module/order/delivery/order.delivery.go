package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/validators"
)

type OrderHandler struct {
	orderUsecase domain.OrderUseCase
	validator    *validator.Validate
	logger       logger.ILogger
}

// InitOrderController will initialize the article's HTTP controller
func InitOrderController(e *echo.Echo, us domain.OrderUseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares) {
	handler := &OrderHandler{
		orderUsecase: us,
		validator:    validator,
		logger:       logger,
	}
	dashboard := e.Group("api/v1/admin/order")
	{
		dashboard.GET("", handler.ListOrderForDashboard)
	}
	mobile := e.Group("api/v1/mobile/order")
	{
		mobile.POST("/calculate-order-cost", handler.CalculateOrderCost)
		dashboard.GET("", handler.ListOrderForMobile)
	}
}

func (a *OrderHandler) ListOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	orders, errResp := a.orderUsecase.ListOrderForDashboard(ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, orders)

}

func (a *OrderHandler) ListOrderForMobile(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	orders, errResp := a.orderUsecase.ListOrderForDashboard(ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, orders)

}

func (a *OrderHandler) CalculateOrderCost(c echo.Context) error {
	ctx := c.Request().Context()

	var calculateOrderCostDto order.CalculateOrderCostDto

	err := c.Bind(&calculateOrderCostDto)
	if err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	validateErr := calculateOrderCostDto.Validate(ctx, a.validator)
	if validateErr.IsError {
		return validators.ErrorStatusUnprocessableEntity(c, validateErr)
	}
	orderCalculate, errResp := a.orderUsecase.CalculateOrderCost(ctx, &calculateOrderCostDto)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		return validators.ErrorResp(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"order_calculate": orderCalculate})
}
