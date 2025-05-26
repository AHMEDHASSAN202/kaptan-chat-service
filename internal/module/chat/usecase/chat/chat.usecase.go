package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"kaptan/internal/module/chat/builder"
	"kaptan/internal/module/chat/consts"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/gate"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"kaptan/pkg/websocket"
)

type ChatUseCase struct {
	repo             domain.ChatRepository
	logger           logger.ILogger
	gate             *gate.Gate
	websocketManager *websocket.ChannelManager
	driverRepo       domain2.DriverRepository
}

func NewChatUseCase(repo domain.ChatRepository, driverRepo domain2.DriverRepository, gate *gate.Gate, websocketManager *websocket.ChannelManager, logger logger.ILogger) domain.ChatUseCase {
	return &ChatUseCase{
		repo:             repo,
		logger:           logger,
		gate:             gate,
		driverRepo:       driverRepo,
		websocketManager: websocketManager,
	}
}

func (u ChatUseCase) GetChats(ctx context.Context, dto *dto.GetChats) (app.ListChatResponse, validators.ErrorResponse) {
	privateChats := u.repo.PrivateChats(ctx, dto)
	chatsResponse := builder.ChatsResponseBuilder(privateChats)
	return chatsResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) GetChatMessages(ctx context.Context, dto *dto.GetChatMessage) (*app.MessagesResponse, validators.ErrorResponse) {
	messages, pagination := u.repo.GetChatMessages(ctx, dto)
	messagesResponse := builder.MessagesResponseBuilder(messages, pagination)
	return messagesResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*app.ChatResponse, validators.ErrorResponse) {
	chat, message, err := u.repo.AddPrivateChat(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	driver, err := u.driverRepo.Find(&ctx, uint(chat.UserId))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	chatResponse := builder.ChatResponseBuilder(chat, driver)

	contentJson, _ := json.Marshal(chatResponse)
	myClient := u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId))
	anotherClient := u.websocketManager.GetClient(utils.GetClientUserId(message.SenderType, fmt.Sprintf("%v", message.SenderId)))

	u.logger.Info("myClient", myClient)
	u.logger.Info("anotherClient", anotherClient)

	u.websocketManager.JoinChannel(myClient, chatResponse.Channel)
	u.websocketManager.JoinChannel(anotherClient, chatResponse.Channel)

	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID: chatResponse.Channel,
		Content:   string(contentJson),
		Action:    consts.START_CHAT_ACTION,
	}

	u.addUnreadMessage(myClient, chatResponse.Channel)

	return chatResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) AcceptPrivateChat(ctx context.Context, dto *dto.AcceptPrivateChat) (*app.ChatResponse, validators.ErrorResponse) {
	chat, err := u.repo.AcceptPrivateChat(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	chatResponse := builder.ChatResponseBuilder(chat, nil)

	contentJson, _ := json.Marshal(chatResponse)

	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID: chatResponse.Channel,
		Content:   string(contentJson),
		Action:    consts.CHANCE_CHAT_STATUS_ACTION,
	}

	return chatResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) RejectOffer(ctx context.Context, dto *dto.RejectOffer) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.RejectOffer(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	messageResponse := builder.MessageResponseBuilder(message)
	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) GetChat(ctx context.Context, dto *dto.GetChat) (*app.ChatResponse, validators.ErrorResponse) {
	chat, err := u.repo.GetChat(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	driver, err := u.driverRepo.Find(&ctx, uint(chat.UserId))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return builder.ChatResponseBuilder(chat, driver), validators.ErrorResponse{}
}

func (u ChatUseCase) SendMessage(ctx context.Context, dto *dto.SendMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.StoreMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	messageResponse := builder.MessageResponseBuilder(message)

	contentJson, _ := json.Marshal(messageResponse)
	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID: messageResponse.Channel,
		Content:   string(contentJson),
		Action:    consts.ADD_MESSAGE_ACTION,
	}

	myClient := u.websocketManager.GetClient(utils.GetClientUserId(dto.CauserType, dto.CauserId))
	u.addUnreadMessage(myClient, messageResponse.Channel)

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.UpdateMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	messageResponse := builder.MessageResponseBuilder(message)

	contentJson, _ := json.Marshal(messageResponse)
	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID: messageResponse.Channel,
		Content:   string(contentJson),
		Action:    consts.UPDATE_MESSAGE_ACTION,
	}

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*app.MessageResponse, validators.ErrorResponse) {
	message, err := u.repo.DeleteMessage(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	messageResponse := builder.MessageResponseBuilder(message)

	contentJson, _ := json.Marshal(messageResponse)
	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID: messageResponse.Channel,
		Content:   string(contentJson),
		Action:    consts.DELETE_MESSAGE_ACTION,
	}

	return messageResponse, validators.ErrorResponse{}
}

func (u ChatUseCase) addUnreadMessage(exceptClient *websocket.Client, channel string) {
	unreadJson, _ := json.Marshal(map[string]string{"channel": channel})
	u.websocketManager.Broadcast <- websocket.Message{
		ChannelID:    channel,
		Content:      string(unreadJson),
		Action:       consts.UNREAD_MESSAGE_ACTION,
		ExceptClient: exceptClient,
	}
}
