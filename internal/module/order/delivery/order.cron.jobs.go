package delivery

import (
	"context"
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

// PickedUpOrdersByCronJob and CancelledOrdersByCronJob
func (a *OrderHandler) CompleteOrdersProcessByCronJob(c echo.Context) error {
	//ctx := c.Request().Context()

	go func() {
		errResp := a.orderUsecase.CronJobPickedOrders(context.Background())
		if errResp.IsError {
			a.logger.Error("CompleteOrdersProcessByCronJob CronJobPickedOrders Err : ", errResp.ErrorMessageObject.Text)
		}
	}()

	go func() {
		errResp := a.orderUsecase.CronJobCancelOrders(context.Background())
		if errResp.IsError {
			a.logger.Error("CompleteOrdersProcessByCronJob CronJobCancelOrders Err : ", errResp.ErrorMessageObject.Text)
		}
	}()

	return validators.SuccessResponse(c, map[string]interface{}{})
}
