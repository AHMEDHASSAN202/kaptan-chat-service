package app

import (
	"gorm.io/gorm"
	"time"
)

type MessageResponse struct {
	ID                      uint           `json:"id"`
	Channel                 string         `json:"channel"`
	SenderType              string         `json:"sender_type"`
	SenderId                int64          `json:"sender_id"`
	TransferId              *int64         `json:"transfer_id"`
	User                    interface{}    `json:"user"`
	BrandId                 *int64         `json:"brand_id"`
	Message                 string         `json:"message"`
	MessageType             string         `json:"message_type"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               *time.Time     `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"deleted_at"`
	CountChannels           int            `json:"count_channels"`
	TransferOffersRequested bool           `json:"transfer_offers_requested"`
	TransferOfferStatus     *string        `json:"transfer_offer_status"`
	Chat                    *ChatResponse  `json:"chat"`
}
