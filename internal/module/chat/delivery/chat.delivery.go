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
		mobile.POST("", handler.AddPrivateChat)
		mobile.GET("", handler.GetChats)
		mobile.PUT("/:id/enable", handler.GetChats)
	}

	{
		mobile.POST("/message", handler.SendMessage)
		mobile.PUT("/message/:id", handler.UpdateMessage)
		mobile.DELETE("/message/:id", handler.DeleteMessage)
	}
}

func (a *ChatHandler) GetChats(c echo.Context) error {
	ctx := c.Request().Context()

	chatsDto := dto.GetChats{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &chatsDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&chatsDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := chatsDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	chats, errResp := a.chatUsecase.GetChats(ctx, &chatsDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"chats": chats})
}

func (a *ChatHandler) AddPrivateChat(c echo.Context) error {
	ctx := c.Request().Context()

	privateChannelDto := dto.AddPrivateChat{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &privateChannelDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&privateChannelDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := privateChannelDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	chat, errResp := a.chatUsecase.AddPrivateChat(ctx, &privateChannelDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"chat": chat})
}

func (a *ChatHandler) EnablePrivateChat(c echo.Context) error {
	ctx := c.Request().Context()

	enablePrivateChatDto := dto.EnablePrivateChat{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &enablePrivateChatDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&enablePrivateChatDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := enablePrivateChatDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	chat, errResp := a.chatUsecase.EnablePrivateChat(ctx, &enablePrivateChatDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"chat": chat})
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

	validationErr := messageDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
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

	validationErr := messageDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
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

	validationErr := messageDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	message, errResp := a.chatUsecase.DeleteMessage(ctx, &messageDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"message": message})
}
