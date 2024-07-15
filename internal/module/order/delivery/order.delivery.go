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
		mobile.POST("", handler.CreateOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/cancel", handler.CancelOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/user-rejection-reason/:status", handler.UserRejectionReason)
	}
}

func (a *OrderHandler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var orderDto order.CreateOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&orderDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	orderResponse, errResp := a.orderUsecase.StoreOrder(ctx, &orderDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
func (a *OrderHandler) CancelOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var orderDto order.CancelOrderDto

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

	orderResponse, errResp := a.orderUsecase.UserCancelOrder(ctx, &orderDto)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
func (a *OrderHandler) UserRejectionReason(c echo.Context) error {
	ctx := c.Request().Context()

	status := c.Param("status")

	rejectionReasons, errResp := a.orderUsecase.UserRejectionReasons(ctx, status, "")
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		return validators.ErrorResp(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"user_rejection_reasons": rejectionReasons})
}

func (a *OrderHandler) ListOrderForDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDto

	_ = c.Bind(&payload)

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
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"order_calculate": orderCalculate})
}
