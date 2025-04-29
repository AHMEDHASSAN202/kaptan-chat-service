package chat

import (
	"context"
	"encoding/json"
	"kaptan/internal/module/chat/builder"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	"kaptan/pkg/gate"
	"kaptan/pkg/logger"
	"kaptan/pkg/validators"
	"kaptan/pkg/websocket"
)

type ChatUseCase struct {
	repo             domain.ChatRepository
	logger           logger.ILogger
	gate             *gate.Gate
	websocketManager *websocket.ChannelManager
}

func NewChatUseCase(repo domain.ChatRepository, gate *gate.Gate, websocketManager *websocket.ChannelManager, logger logger.ILogger) domain.ChatUseCase {
	return &ChatUseCase{
		repo:             repo,
		logger:           logger,
		gate:             gate,
		websocketManager: websocketManager,
	}
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
		Action:    "new-message",
	}

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
		Action:    "update-message",
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
		Action:    "delete-message",
	}

	return messageResponse, validators.ErrorResponse{}
}
