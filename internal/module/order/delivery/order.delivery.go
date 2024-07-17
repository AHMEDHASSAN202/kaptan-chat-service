package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	"samm/pkg/middlewares/kitchen"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/validators"
)

type OrderHandler struct {
	orderUsecase domain.OrderUseCase
	validator    *validator.Validate
	logger       logger.ILogger
}

// InitOrderController will initialize the article's HTTP controller
func InitOrderController(e *echo.Echo, us domain.OrderUseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares, kitchenMiddlewares *kitchen.ProviderMiddlewares) {
	handler := &OrderHandler{
		orderUsecase: us,
		validator:    validator,
		logger:       logger,
	}
	dashboard := e.Group("api/v1/admin/order")
	{
		dashboard.GET("", handler.ListOrderForDashboard)
		dashboard.GET("/:id", handler.FindOrderForDashboard)
	}
	mobile := e.Group("api/v1/mobile/order")
	{
		mobile.POST("/calculate-order-cost", handler.CalculateOrderCost)
		mobile.POST("", handler.CreateOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("", handler.ListOrderForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/:id", handler.FindOrderForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/toggle-favourite", handler.ToggleOrderFavourite, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/cancel", handler.CancelOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/arrived", handler.ArrivedOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/user-rejection-reason/:status", handler.UserRejectionReason)
	}

	kitchen := e.Group("api/v1/kitchen/order")
	{
		kitchen.PUT("/:id/accept", handler.KitchenToAccept, kitchenMiddlewares.AuthMiddleware)
		kitchen.PUT("/:id/rejected", handler.KitchenToRejected, kitchenMiddlewares.AuthMiddleware)
		kitchen.GET("/rejection-reason/:status", handler.KitchenRejectionReason, kitchenMiddlewares.AuthMiddleware)
		kitchen.PUT("/:id/ready-for-pickup", handler.KitchenToReadyForPickup, kitchenMiddlewares.AuthMiddleware)
		kitchen.PUT("/:id/picked-up", handler.KitchenToPickedUp, kitchenMiddlewares.AuthMiddleware)
		kitchen.PUT("/:id/no-show", handler.KitchenToNoShow, kitchenMiddlewares.AuthMiddleware)
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
func (a *OrderHandler) ArrivedOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var orderDto order.ArrivedOrderDto

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

	orderResponse, errResp := a.orderUsecase.UserArrivedOrder(ctx, &orderDto)
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

func (a *OrderHandler) FindOrderForMobile(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.FindOrderMobileDto

	orderId := c.Param("id")
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	payload.OrderId = orderId

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	order, errResp := a.orderUsecase.FindOrderForMobile(&ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, order)
}

func (a *OrderHandler) ListOrderForMobile(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ListOrderDtoForMobile

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	payload.Pagination.SetDefault()

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	orders, errResp := a.orderUsecase.ListOrderForMobile(ctx, &payload)
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

func (a *OrderHandler) ToggleOrderFavourite(c echo.Context) error {
	ctx := c.Request().Context()
	var payload order.ToggleOrderFavDto

	orderId := c.Param("id")
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	payload.OrderId = orderId

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	errResp := a.orderUsecase.ToggleOrderFavourite(&ctx, payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
