package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/payment/domain"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type PaymentHandler struct {
	paymentUseCase domain.PaymentUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitPaymentController(e *echo.Echo, us domain.PaymentUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &PaymentHandler{
		paymentUseCase: us,
		validator:      validator,
		logger:         logger,
	}
	mobile := e.Group("api/v1/mobile/payment")
	mobile.POST("", handler.pay)
}

func (a *PaymentHandler) pay(c echo.Context) error {
	//ctx := c.Request().Context()

	//var payload user.CreateUserDto
	//err := c.Bind(&payload)
	//if err != nil {
	//	return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	//}
	//
	//validationErr := payload.Validate(c, a.validator)
	//if validationErr.IsError {
	//	a.logger.Error(validationErr)
	//	return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	//}
	//
	//errResp := a.userUsecase.StoreUser(ctx, &payload)
	//if errResp.IsError {
	//	a.logger.Error(errResp)
	//	return validators.ErrorStatusBadRequest(c, errResp)
	//}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
