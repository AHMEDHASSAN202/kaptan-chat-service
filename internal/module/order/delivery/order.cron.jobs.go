package delivery

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"samm/pkg/validators"
)

func (a *OrderHandler) TimedOutOrdersByCronJob(c echo.Context) error {
	ctx := c.Request().Context()

	errResp := a.orderUsecase.CronJobTimedOutOrders(ctx)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		errResp.StatusCode = http.StatusBadRequest
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *OrderHandler) PickedUpOrdersByCronJob(c echo.Context) error {
	ctx := c.Request().Context()

	errResp := a.orderUsecase.CronJobPickedOrders(ctx)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		errResp.StatusCode = http.StatusBadRequest
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *OrderHandler) CancelledOrdersByCronJob(c echo.Context) error {
	ctx := c.Request().Context()

	errResp := a.orderUsecase.CronJobCancelOrders(ctx)
	if errResp.IsError {
		a.logger.Error(errResp.ErrorMessageObject.Text)
		errResp.StatusCode = http.StatusBadRequest
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
