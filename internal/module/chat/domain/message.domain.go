package domain

import (
	"gorm.io/gorm"
	"kaptan/pkg/database/mysql/custom_types"
)

type Message struct {
	gorm.Model
	Channel                 string               `gorm:"column:channel;index:channel"`
	SenderType              string               `gorm:"column:sender_type;index:sender_entity"`
	SenderId                int64                `gorm:"column:sender_id;index:sender_entity"`
	TransferId              *int64               `gorm:"column:transfer_id"`
	BrandId                 *int64               `gorm:"column:brand_id"`
	Message                 string               `gorm:"column:message"`
	MessageType             string               `gorm:"column:message_type"` //test, image or file
	User                    custom_types.JSONMap `gorm:"column:user;type:json"`
	IsPrivate               bool                 `gorm:"column:is_private"`
	CountChannels           int                  `gorm:"column:count_channels;default:0"`
	TransferOffersRequested bool                 `gorm:"column:transfer_offers_requested;default:0"`
	TransferOfferStatus     *string              `gorm:"column:transfer_offer_status;default:null"`
}

type MessageUseCase interface {
}

type MessageRepository interface {
}
