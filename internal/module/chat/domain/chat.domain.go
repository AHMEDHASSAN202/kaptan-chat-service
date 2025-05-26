package domain

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	"kaptan/pkg/database/mysql"
	"kaptan/pkg/database/mysql/custom_types"
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
	UnreadMessagesCount int                  `gorm:"column:unread_messages_count;default:0"`
	Status              string               `gorm:"column:status"`
	OpenedBy            uint                 `gorm:"column:opened_by;index:opened_by_message_id_index"`
}

type ChatUseCase interface {
	GetChats(ctx context.Context, dto *dto.GetChats) (app.ListChatResponse, validators.ErrorResponse)
	GetChat(ctx context.Context, dto *dto.GetChat) (*app.ChatResponse, validators.ErrorResponse)
	GetChatMessages(ctx context.Context, dto *dto.GetChatMessage) (*app.MessagesResponse, validators.ErrorResponse)
	AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*app.ChatResponse, validators.ErrorResponse)
	AcceptPrivateChat(ctx context.Context, dto *dto.AcceptPrivateChat) (*app.ChatResponse, validators.ErrorResponse)
	SendMessage(ctx context.Context, message *dto.SendMessage) (*app.MessageResponse, validators.ErrorResponse)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*app.MessageResponse, validators.ErrorResponse)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*app.MessageResponse, validators.ErrorResponse)
	RejectOffer(ctx context.Context, dto *dto.RejectOffer) (*app.MessageResponse, validators.ErrorResponse)
}

type ChatRepository interface {
	PrivateChats(ctx context.Context, dto *dto.GetChats) []*Chat
	GetChat(ctx context.Context, dto *dto.GetChat) (*Chat, error)
	GetChatMessages(ctx context.Context, dto *dto.GetChatMessage) ([]*Message, *mysql.Pagination)
	AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*Chat, *Message, error)
	AcceptPrivateChat(ctx context.Context, dto *dto.AcceptPrivateChat) (*Chat, error)
	StoreMessage(ctx context.Context, message *dto.SendMessage) (*Message, error)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*Message, error)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*Message, error)
	RejectOffer(ctx context.Context, dto *dto.RejectOffer) (*Message, error)
}
