package builder

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
	"kaptan/pkg/database/mysql"
)

func MessagesResponseBuilder(messages []*domain.Message, pagination *mysql.Pagination) *app.MessagesResponse {
	messagesResponse := &app.MessagesResponse{
		Docs: make([]*app.MessageResponse, 0),
		Meta: pagination,
	}
	for _, message := range messages {
		messagesResponse.Docs = append(messagesResponse.Docs, MessageResponseBuilder(message))
	}
	return messagesResponse
}
