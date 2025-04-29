package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/pkg/logger"
	commmon "kaptan/pkg/middlewares/common"
	usermiddleware "kaptan/pkg/middlewares/user"
	"kaptan/pkg/validators"
)

type ChatHandler struct {
	chatUsecase domain.ChatUseCase
	validator   *validator.Validate
	logger      logger.ILogger
}

// InitChatController will initialize the article's HTTP controller
func InitChatController(e *echo.Echo, us domain.ChatUseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &ChatHandler{
		chatUsecase: us,
		validator:   validator,
		logger:      logger,
	}

	mobile := e.Group("api/v1/app/chat")
	mobile.Use(userMiddleware.AuthenticationMiddleware("driver"))
	{
		mobile.POST("/message", handler.SendMessage)
		mobile.PUT("/message/:id", handler.UpdateMessage)
	}
}

func (a *ChatHandler) SendMessage(c echo.Context) error {
	ctx := c.Request().Context()

	messageDto := dto.SendMessage{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	message, errResp := a.chatUsecase.SendMessage(ctx, &messageDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"message": message})
}

func (a *ChatHandler) UpdateMessage(c echo.Context) error {
	ctx := c.Request().Context()

	messageDto := dto.UpdateMessage{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	message, errResp := a.chatUsecase.UpdateMessage(ctx, &messageDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"message": message})
}

func (a *ChatHandler) DeleteMessage(c echo.Context) error {
	ctx := c.Request().Context()

	messageDto := dto.DeleteMessage{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&messageDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	message, errResp := a.chatUsecase.DeleteMessage(ctx, &messageDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"message": message})
}
