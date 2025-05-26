package builder

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/responses/app"
)

func ChatsResponseBuilder(chats []*domain.Chat) app.ListChatResponse {
	response := app.ListChatResponse{}
	for _, chat := range chats {
		response = append(response, ChatResponseBuilder(chat, nil))
	}
	return response
}
