package builder

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
)

func MessageResponseBuilder(message *domain.Message) *app.MessageResponse {
	messageResponse := &app.MessageResponse{
		ID:          message.ID,
		Channel:     message.Channel,
		SenderType:  message.SenderType,
		SenderId:    message.SenderId,
		Message:     message.Message,
		MessageType: message.MessageType,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   &message.UpdatedAt,
		DeletedAt:   message.DeletedAt,
	}

	if message.TransferId.Valid {
		messageResponse.TransferId = &message.TransferId.Int64
	}

	if message.OwnerTransferId.Valid {
		messageResponse.OwnerTransferId = &message.OwnerTransferId.Int64
	}

	if message.BrandId.Valid {
		messageResponse.BrandId = &message.BrandId.Int64
	}

	if message.CreatedAt == message.UpdatedAt {
		messageResponse.UpdatedAt = nil
	}

	return messageResponse
}
