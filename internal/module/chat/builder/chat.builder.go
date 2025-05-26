package builder

import (
	"fmt"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/utils"
)

func ChatResponseBuilder(chat *domain.Chat, driver *domain2.Driver) *app.ChatResponse {
	response := &app.ChatResponse{
		ID:                  chat.ID,
		Channel:             chat.Channel,
		Name:                getChatName(chat),
		CreatedAt:           chat.CreatedAt,
		UpdatedAt:           &chat.UpdatedAt,
		DeletedAt:           chat.DeletedAt,
		IsOwner:             chat.IsOwner,
		User:                chat.User,
		TransferId:          chat.TransferId,
		LastMessage:         chat.LastMessage,
		Status:              chat.Status,
		UnreadMessagesCount: chat.UnreadMessagesCount,
	}
	if driver != nil {
		response.User = utils.StructToMap(driver.ToResponse(), "json")
	}
	if chat.CreatedAt == chat.UpdatedAt {
		response.UpdatedAt = nil
	}
	return response
}

func getChatName(chat *domain.Chat) string {
	if val, ok := chat.User["name"]; ok {
		return fmt.Sprintf("%s #%d", val, chat.ID)
	}
	return chat.Channel
}
