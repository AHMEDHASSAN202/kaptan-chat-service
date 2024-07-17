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

	orderResponse, errResp := a.orderUsecase.KitchenAcceptOrder(ctx, &acceptDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}

func (a *OrderHandler) KitchenToRejected(c echo.Context) error {
	ctx := c.Request().Context()

	var dto kitchen.RejectedOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	orderResponse, errResp := a.orderUsecase.KitchenRejectedOrder(ctx, &dto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}

func (a *OrderHandler) KitchenRejectionReason(c echo.Context) error {
	ctx := c.Request().Context()

	status := c.Param("status")

	rejectionReasons, errResp := a.orderUsecase.KitchenRejectionReasons(ctx, status, "")
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		return validators.ErrorResp(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"rejection_reasons": rejectionReasons})
}

func (a *OrderHandler) KitchenToReadyForPickup(c echo.Context) error {
	ctx := c.Request().Context()

	var dto kitchen.ReadyForPickupOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	orderResponse, errResp := a.orderUsecase.KitchenReadyForPickupOrder(ctx, &dto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}

func (a *OrderHandler) KitchenToPickedUp(c echo.Context) error {
	ctx := c.Request().Context()

	var dto kitchen.PickedUpOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	orderResponse, errResp := a.orderUsecase.KitchenPickedUpOrder(ctx, &dto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}

func (a *OrderHandler) KitchenToNoShow(c echo.Context) error {
	ctx := c.Request().Context()

	var dto kitchen.NoShowOrderDto

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&dto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	orderResponse, errResp := a.orderUsecase.KitchenNoShowOrder(ctx, &dto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"order": orderResponse})
}
