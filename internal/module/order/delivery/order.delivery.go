package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
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
func InitOrderController(e *echo.Echo, us domain.OrderUseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares, kitchenMiddlewares *kitchen.ProviderMiddlewares) {
	handler := &OrderHandler{
		orderUsecase: us,
		validator:    validator,
		logger:       logger,
	}
	dashboard := e.Group("api/v1/admin/order")
	dashboard.Use(adminMiddlewares.AuthMiddleware)
	{
		dashboard.GET("", handler.ListOrderForDashboard, commonMiddlewares.PermissionMiddleware("list-orders"))
		dashboard.GET("/:id", handler.FindOrderForDashboard, commonMiddlewares.PermissionMiddleware("find-order"))
		dashboard.PUT("/:id/cancel", handler.CancelOrderForDashboard, commonMiddlewares.PermissionMiddleware("update-order-status"))
		dashboard.PUT("/:id/pickedup", handler.PickedUpOrderForDashboard, commonMiddlewares.PermissionMiddleware("update-order-status"))
	}
	mobile := e.Group("api/v1/mobile/order")
	{
		mobile.POST("/calculate-order-cost", handler.CalculateOrderCost)
		mobile.POST("", handler.CreateOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/in-progress", handler.ListInprogressOrdersForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/completed", handler.ListCompletedOrdersForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/last", handler.ListLastOrdersForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/:id", handler.FindOrderForMobile, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/toggle-favourite", handler.ToggleOrderFavourite, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/cancel", handler.CancelOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.PUT("/:id/arrived", handler.ArrivedOrder, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
		mobile.GET("/user-rejection-reason/:status", handler.UserRejectionReason)
		mobile.PUT("/:id/report-missing-item", handler.ReportMissingItem, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
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

	cronJobs := e.Group("api/v1/cron-jobs/order")
	{
		cronJobs.POST("/time-out", handler.TimedOutOrdersByCronJob)
		cronJobs.POST("/picked-up", handler.PickedUpOrdersByCronJob)
		cronJobs.POST("/cancel", handler.CancelledOrdersByCronJob)
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

func (a *OrderHandler) ListInprogressOrdersForMobile(c echo.Context) error {
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

	orders, errResp := a.orderUsecase.ListInprogressOrdersForMobile(ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, orders)

}

func (a *OrderHandler) ListCompletedOrdersForMobile(c echo.Context) error {
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

	orders, errResp := a.orderUsecase.ListCompletedOrdersForMobile(ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, orders)

}

func (a *OrderHandler) ListLastOrdersForMobile(c echo.Context) error {
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

	orders, errResp := a.orderUsecase.ListLastOrdersForMobile(ctx, &payload)
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

func (a *OrderHandler) ReportMissingItem(c echo.Context) error {
	ctx := c.Request().Context()

	var payload order.ReportMissingItemDto
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := binder.BindPathParams(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := binder.BindBody(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	_, errResp := a.orderUsecase.ReportMissedItem(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
