package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
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
	mobile := e.Group("api/v1/mobile")
	mobile.POST("/:transactionType/pay", handler.pay)
	myfatoorah := e.Group("api/v1/myfatoorah")
	myfatoorah.POST("/webhook", handler.MyfatoorahWebhook)
}

func (a *PaymentHandler) pay(c echo.Context) error {
	ctx := c.Request().Context()

	var payload payment.PayDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	transactionType := c.Param("transactionType")
	payload.TransactionType = transactionType
	a.logger.Info(payload)
	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	//
	res, errResp := a.paymentUseCase.Pay(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}
func (a *PaymentHandler) MyfatoorahWebhook(c echo.Context) error {
	ctx := c.Request().Context()

	var payload payment.MyFatoorahWebhookPayload
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	a.logger.Info(payload)
	//
	res, errResp := a.paymentUseCase.Pay(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}
