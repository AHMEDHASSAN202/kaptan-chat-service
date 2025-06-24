package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/dto"
	"kaptan/pkg/logger"
	commmon "kaptan/pkg/middlewares/common"
	usermiddleware "kaptan/pkg/middlewares/user"
	"kaptan/pkg/validators"
)

type TransferHandler struct {
	transferUseCase domain.UseCase
	validator       *validator.Validate
	logger          logger.ILogger
}

// InitChatController will initialize the article's HTTP controller
func InitTransferController(e *echo.Echo, us domain.UseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &TransferHandler{
		transferUseCase: us,
		validator:       validator,
		logger:          logger,
	}

	mobile := e.Group("api/v1/app/transfer")
	mobile.Use(userMiddleware.AuthenticationMiddleware("driver"))
	{
		mobile.PUT("/:id/start", handler.StartTransfer)
		mobile.PUT("/:id/end", handler.EndTransfer)
	}
}

func (a *TransferHandler) StartTransfer(c echo.Context) error {
	ctx := c.Request().Context()

	startDto := dto.StartTransfer{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &startDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&startDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := startDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.transferUseCase.StartTransfer(&ctx, &startDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *TransferHandler) EndTransfer(c echo.Context) error {
	ctx := c.Request().Context()

	endDto := dto.EndTransfer{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &endDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&endDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := endDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.transferUseCase.EndTransfer(&ctx, &endDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
