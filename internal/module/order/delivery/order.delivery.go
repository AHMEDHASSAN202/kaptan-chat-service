package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type OrderHandler struct {
	orderUsecase domain.OrderUseCase
	validator    *validator.Validate
	logger       logger.ILogger
}

// InitOrderController will initialize the article's HTTP controller
func InitOrderController(e *echo.Echo, us domain.OrderUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &OrderHandler{
		orderUsecase: us,
		validator:    validator,
		logger:       logger,
	}
	dashboard := e.Group("api/v1/admin/order")
	{
		dashboard.GET("", handler.ListOrder)
		dashboard.GET("/:id", handler.FindOrder)
	}
	mobile := e.Group("api/v1/mobile/order")
	{
		mobile.POST("/calculate-order-cost", handler.CalculateOrderCost)
		mobile.GET("/:id", handler.FindOrder)
	}
}

func (a *OrderHandler) FindOrder(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.orderUsecase.FindOrder(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"order": data})
}

func (a *OrderHandler) ListOrder(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, paginationResult, errResp := a.orderUsecase.ListOrder(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
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
		a.logger.Error(validateErr.ErrorMessageObject.Text)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	orderCalculate, errResp := a.orderUsecase.CalculateOrderCost(ctx, &calculateOrderCostDto)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"order_calculate": orderCalculate})
}
