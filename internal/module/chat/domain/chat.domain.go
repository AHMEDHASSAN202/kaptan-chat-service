package domain

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain/custom_types"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	"kaptan/pkg/validators"
)

type Chat struct {
	gorm.Model
	Channel             string               `gorm:"column:channel;index:channel"`
	UserType            string               `gorm:"column:user_type;index:user_entity"`
	UserId              int                  `gorm:"column:user_id;index:user_entity"`
	TransferId          *int64               `gorm:"column:transfer_id"`
	IsOwner             bool                 `gorm:"column:is_owner"`
	User                custom_types.JSONMap `gorm:"column:user;type:json"`
	LastMessage         custom_types.JSONMap `gorm:"column:last_message;type:json"`
	UnreadMessagesCount int                  `gorm:"column:unread_messages_count"`
	Disabled            bool                 `gorm:"column:disabled"`
}

type ChatUseCase interface {
	GetChats(ctx context.Context, dto *dto.GetChats) (app.ListChatResponse, validators.ErrorResponse)
	AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*app.ChatResponse, validators.ErrorResponse)
	EnablePrivateChat(ctx context.Context, dto *dto.EnablePrivateChat) (*app.ChatResponse, validators.ErrorResponse)
	SendMessage(ctx context.Context, message *dto.SendMessage) (*app.MessageResponse, validators.ErrorResponse)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*app.MessageResponse, validators.ErrorResponse)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*app.MessageResponse, validators.ErrorResponse)
}

type ChatRepository interface {
	PrivateChats(ctx context.Context, dto *dto.GetChats) []*Chat
	AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*Chat, *Message, error)
	EnablePrivateChat(ctx context.Context, dto *dto.EnablePrivateChat) (*Chat, error)
	StoreMessage(ctx context.Context, message *dto.SendMessage) (*Message, error)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*Message, error)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*Message, error)
}
