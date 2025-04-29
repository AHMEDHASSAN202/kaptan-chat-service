package migrations

import (
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain"
)

func NewChatChannelsMigration(db *gorm.DB) {
	db.AutoMigrate(domain.Chat{})
}
