package migrations

import (
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain"
)

func NewMessagesMigration(db *gorm.DB) {
	db.AutoMigrate(domain.Message{})
}
