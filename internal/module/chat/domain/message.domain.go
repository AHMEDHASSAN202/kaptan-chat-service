package domain

import (
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain/custom_types"
)

type Message struct {
	gorm.Model
	Channel     string               `gorm:"column:channel;index:channel"`
	SenderType  string               `gorm:"column:sender_type;index:sender_entity"`
	SenderId    int64                `gorm:"column:sender_id;index:sender_entity"`
	TransferId  *int64               `gorm:"column:transfer_id"`
	BrandId     *int64               `gorm:"column:brand_id"`
	Message     string               `gorm:"column:message"`
	MessageType string               `gorm:"column:message_type"` //test, image or file
	User        custom_types.JSONMap `gorm:"column:user;type:json"`
}

type MessageUseCase interface {
}

type MessageRepository interface {
}
