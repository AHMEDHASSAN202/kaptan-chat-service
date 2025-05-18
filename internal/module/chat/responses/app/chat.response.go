package app

import (
	"gorm.io/gorm"
	"time"
)

type ChatResponse struct {
	ID                  uint                   `json:"id"`
	Channel             string                 `json:"channel"`
	Name                string                 `json:"name"`
	TransferId          *int64                 `json:"transfer_id"`
	IsOwner             bool                   `json:"is_owner"`
	UnreadMessagesCount int                    `json:"unread_messages_count"`
	User                map[string]interface{} `json:"user"`
	LastMessage         map[string]interface{} `json:"last_message"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           *time.Time             `json:"updated_at"`
	DeletedAt           gorm.DeletedAt         `json:"deleted_at"`
	Disabled            bool                   `json:"disabled"`
}
