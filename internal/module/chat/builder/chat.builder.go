package builder

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
)

func ChatResponseBuilder(chat *domain.Chat) *app.ChatResponse {
	response := &app.ChatResponse{
		ID:          chat.ID,
		Channel:     chat.Channel,
		CreatedAt:   chat.CreatedAt,
		UpdatedAt:   &chat.UpdatedAt,
		DeletedAt:   chat.DeletedAt,
		IsOwner:     chat.IsOwner,
		User:        chat.User,
		TransferId:  chat.TransferId,
		LastMessage: chat.LastMessage,
		Disabled:    chat.Disabled,
	}
	if chat.CreatedAt == chat.UpdatedAt {
		response.UpdatedAt = nil
	}
	return response
}
