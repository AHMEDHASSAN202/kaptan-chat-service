package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/pkg/database/redis"
	"samm/pkg/logger"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/validators"
)

type PaymentHandler struct {
	paymentUseCase domain.PaymentUseCase
	validator      *validator.Validate
	logger         logger.ILogger
	redis          *redis.RedisClient
	userMiddleware *usermiddleware.Middlewares
}

// InitUserController will initialize the article's HTTP controller
func InitPaymentController(e *echo.Echo, us domain.PaymentUseCase, validator *validator.Validate, logger logger.ILogger, rdb *redis.RedisClient, userMiddleware *usermiddleware.Middlewares) {
	handler := &PaymentHandler{
		paymentUseCase: us,
		validator:      validator,
		logger:         logger,
		redis:          rdb,
		userMiddleware: userMiddleware,
	}
	mobile := e.Group("api/v1/mobile")
	mobile.POST("/:transactionType/pay", handler.pay, userMiddleware.AuthenticationMiddleware(false))
	mobile.GET("/payment/:transactionType/:id/status", handler.GetPaymentStatus, userMiddleware.AuthenticationMiddleware(false))

	mobile.PUT("/myfatoorah/update-session", handler.UpdateSession)

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
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	userId := payload.CauserId

	//duration := time.Now().UTC().Add(10 * time.Second).Sub(time.Now().UTC())
	//lock, er := a.redis.Lock("USER_PAYMENT_"+userId, userId, duration)
	//
	//if er != nil {
	//	return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(er))
	//}
	//if !lock {
	//	return validators.ErrorResp(c, validators.GetErrorResponse(&ctx, localization.PAYMENT_PROCESS_RUNNING, nil, utils.GetAsPointer(http.StatusBadRequest)))
	//}
	//defer a.redis.Unlock("USER_PAYMENT_" + userId)

	transactionType := c.Param("transactionType")
	payload.TransactionType = transactionType
	payload.UserId = userId
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
func (a *PaymentHandler) GetPaymentStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var payload payment.GetPaymentStatus
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	userId := payload.CauserId

	TransactionId := c.Param("id")
	transactionType := c.Param("transactionType")
	payload.TransactionId = TransactionId
	payload.TransactionType = transactionType
	payload.UserId = userId
	//
	res, errResp := a.paymentUseCase.GetPaymentStatus(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}
func (a *PaymentHandler) UpdateSession(c echo.Context) error {
	ctx := c.Request().Context()

	var payload payment.UpdateSession
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	//
	res, errResp := a.paymentUseCase.UpdateSession(ctx, &payload)
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
	res, errResp := a.paymentUseCase.MyFatoorahWebhook(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}
