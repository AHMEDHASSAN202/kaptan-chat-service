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
	dashboard.POST("", handler.StoreOrder)
	dashboard.GET("", handler.ListOrder)
	dashboard.PUT("/:id", handler.UpdateOrder)
	dashboard.GET("/:id", handler.FindOrder)
	dashboard.DELETE("/:id", handler.DeleteOrder)
}
func (a *OrderHandler) StoreOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var payload order.StoreOrderDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.orderUsecase.StoreOrder(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *OrderHandler) UpdateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var payload order.UpdateOrderDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	id := c.Param("id")
	errResp := a.orderUsecase.UpdateOrder(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
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

func (a *OrderHandler) DeleteOrder(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.orderUsecase.DeleteOrder(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
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
