package builder

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
)

func MessageResponseBuilder(message *domain.Message) *app.MessageResponse {
	messageResponse := &app.MessageResponse{
		ID:                      message.ID,
		Channel:                 message.Channel,
		SenderType:              message.SenderType,
		SenderId:                message.SenderId,
		Message:                 message.Message,
		MessageType:             message.MessageType,
		User:                    message.User,
		TransferId:              message.TransferId,
		BrandId:                 message.BrandId,
		CreatedAt:               message.CreatedAt,
		UpdatedAt:               &message.UpdatedAt,
		DeletedAt:               message.DeletedAt,
		CountChannels:           message.CountChannels,
		TransferOffersRequested: message.TransferOffersRequested,
		TransferOfferStatus:     message.TransferOfferStatus,
	}

	if message.CreatedAt == message.UpdatedAt {
		messageResponse.UpdatedAt = nil
	}

	return messageResponse
}
