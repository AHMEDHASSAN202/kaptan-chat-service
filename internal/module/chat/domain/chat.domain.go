package domain

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/dto"
	"kaptan/internal/module/chat/responses/app"
	"kaptan/pkg/validators"
)

type Chat struct {
	gorm.Model
	Channel             string        `gorm:"column:channel;index:channel"`
	UserType            string        `gorm:"column:user_type;index:user_entity"`
	UserId              int           `gorm:"column:user_id;index:user_entity"`
	TransferId          sql.NullInt64 `gorm:"column:transfer_id"`
	OwnerTransferId     sql.NullInt64 `gorm:"column:owner_transfer_id"`
	UnreadMessagesCount bool          `gorm:"column:unread_messages_count"`
}

type ChatUseCase interface {
	SendMessage(ctx context.Context, message *dto.SendMessage) (*app.MessageResponse, validators.ErrorResponse)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*app.MessageResponse, validators.ErrorResponse)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*app.MessageResponse, validators.ErrorResponse)
}

type ChatRepository interface {
	StoreMessage(ctx context.Context, message *dto.SendMessage) (*Message, error)
	UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*Message, error)
	DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*Message, error)
}
