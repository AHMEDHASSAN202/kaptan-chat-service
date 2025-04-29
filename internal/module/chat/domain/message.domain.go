package domain

import (
	"database/sql"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Channel         string        `gorm:"column:channel;index:channel"`
	SenderType      string        `gorm:"column:sender_type;index:sender_entity"`
	SenderId        int64         `gorm:"column:sender_id;index:sender_entity"`
	TransferId      sql.NullInt64 `gorm:"column:transfer_id"`
	OwnerTransferId sql.NullInt64 `gorm:"column:owner_transfer_id"`
	BrandId         sql.NullInt64 `gorm:"column:brand_id"`
	Message         string        `gorm:"column:message"`
	MessageType     string        `gorm:"column:message_type"` //test, image or file
}

type MessageUseCase interface {
}

type MessageRepository interface {
}
