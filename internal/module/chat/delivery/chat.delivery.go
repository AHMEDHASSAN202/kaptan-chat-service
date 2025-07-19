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
		mobile.GET("/:channel", handler.GetChat)
		mobile.GET("/:channel/messages", handler.GetChatMessage)
		mobile.PUT("/:channel/sale-transfer", handler.SaleTransferChat)
	}

	{
		mobile.POST("/message", handler.SendMessage)
		mobile.PUT("/message/:id", handler.UpdateMessage)
		mobile.DELETE("/message/:id", handler.DeleteMessage)
		mobile.PUT("/message/:id/reject", handler.RejectOffer)
		mobile.GET("/unread-messages", handler.UnreadMessages)
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

func (a *ChatHandler) GetChat(c echo.Context) error {
	ctx := c.Request().Context()

	chatsDto := dto.GetChat{}

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

	chat, transfer, errResp := a.chatUsecase.GetChat(ctx, &chatsDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"chat": chat, "transfer": transfer})
}

func (a *ChatHandler) SaleTransferChat(c echo.Context) error {
	ctx := c.Request().Context()

	chatsDto := dto.SaleTransferChat{}

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

	chat, errResp := a.chatUsecase.SaleTransferChat(ctx, &chatsDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"chat": chat})
}

func (a *ChatHandler) GetChatMessage(c echo.Context) error {
	ctx := c.Request().Context()

	messagesDto := dto.GetChatMessage{}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &messagesDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if err := c.Bind(&messagesDto); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := messagesDto.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	messages, errResp := a.chatUsecase.GetChatMessages(ctx, &messagesDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"messages": messages})
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

func (a *ChatHandler) RejectOffer(c echo.Context) error {
	ctx := c.Request().Context()

	messageDto := dto.RejectOffer{}

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

	message, errResp := a.chatUsecase.RejectOffer(ctx, &messageDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"message": message})
}

func (a *ChatHandler) UnreadMessages(c echo.Context) error {
	ctx := c.Request().Context()

	chatsDto := dto.UnreadMessages{}

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

	count, errResp := a.chatUsecase.UnreadMessages(ctx, &chatsDto)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"unread_messages_count": count})
}
