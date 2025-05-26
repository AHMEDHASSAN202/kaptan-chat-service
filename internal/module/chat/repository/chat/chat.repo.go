package chat

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/consts"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/database/mysql"
	"kaptan/pkg/gate"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"math/rand"
	"strings"
	"time"
)

type ChatRepository struct {
	db               *gorm.DB
	logger           logger.ILogger
	gate             *gate.Gate
	driverRepository domain2.DriverRepository
}

func NewChatRepository(logger logger.ILogger, db *gorm.DB, gate *gate.Gate, driverRepository domain2.DriverRepository) domain.ChatRepository {
	return &ChatRepository{
		db:               db,
		logger:           logger,
		gate:             gate,
		driverRepository: driverRepository,
	}
}

func (r ChatRepository) PrivateChats(ctx context.Context, dto *dto.GetChats) []*domain.Chat {
	var chats []*domain.Chat
	q := r.db.Where("user_type = ?", dto.CauserType).Where("user_id = ?", dto.CauserId)
	if dto.MessageId != nil {
		q.Where("opened_by = ?", dto.MessageId)
	}
	q.Order("updated_at desc").Find(&chats)
	return chats
}

func (r ChatRepository) GetChatMessages(ctx context.Context, dto *dto.GetChatMessage) ([]*domain.Message, *mysql.Pagination) {
	pagination := mysql.Pagination{}
	var messages []*domain.Message
	query := r.db.Model(domain.Message{})
	if strings.ToLower(dto.Channel) != "all" {
		query = query.Where("channel = ?", dto.Channel)
	}
	if strings.ToLower(dto.MyMessage) == "true" {
		query = query.Where("sender_id = ?", dto.CauserId).Where("is_private = ?", 0)
	}
	query.Scopes(mysql.Paginate(&pagination, query, dto.Pagination)).Find(&messages)
	return messages, &pagination
}

func (r ChatRepository) AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*domain.Chat, *domain.Message, error) {
	message := domain.Message{}
	message.ID = dto.MessageId
	r.db.First(&message)

	go func() {
		if err := r.db.Model(&domain.Message{}).Where("id = ?", message.ID).
			UpdateColumn("count_channels", gorm.Expr("count_channels + ?", 1)).Error; err != nil {
			r.logger.Error(err)
		}
	}()

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(999999999999-111111111111+1) + 111111111111
	channel := fmt.Sprintf("private-%v-%v-%d", dto.MessageId, dto.CauserId, randomNumber)

	currentChat := domain.Chat{}
	r.db.Model(&domain.Chat{}).Where("channel = ?", channel).First(&currentChat)
	if currentChat.ID != 0 {
		return &currentChat, &message, nil
	}

	user, err := r.driverRepository.Find(&ctx, uint(*utils.StringToUint(dto.CauserId)))
	if err != nil {
		return nil, nil, err
	}

	chat := &domain.Chat{
		Channel:             channel,
		UserType:            dto.CauserType,
		UserId:              int(*utils.StringToUint(dto.CauserId)),
		TransferId:          message.TransferId,
		UnreadMessagesCount: 1,
		IsOwner:             false,
		LastMessage: map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"brand_id":     message.BrandId,
			"transfer_id":  message.TransferId,
			"message":      message.Message,
			"message_type": message.MessageType,
		},
		Status:   consts.PENDING_CHAT_STATUS,
		User:     message.User,
		OpenedBy: dto.MessageId,
	}

	ownerUserChat := &domain.Chat{
		Channel:             channel,
		UserType:            message.SenderType,
		UserId:              int(message.SenderId),
		TransferId:          message.TransferId,
		UnreadMessagesCount: 1,
		IsOwner:             true,
		LastMessage: map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"brand_id":     message.BrandId,
			"transfer_id":  message.TransferId,
			"message":      message.Message,
			"message_type": message.MessageType,
		},
		Status:   consts.ACCEPT_CHAT_STATUS,
		User:     utils.StructToMap(user.ToResponse(), "json"),
		OpenedBy: dto.MessageId,
	}

	result := r.db.Create([]*domain.Chat{
		chat, ownerUserChat,
	})

	return chat, &message, result.Error
}

func (r ChatRepository) AcceptPrivateChat(ctx context.Context, dto *dto.AcceptPrivateChat) (*domain.Chat, error) {
	var chat *domain.Chat
	r.db.Model(&domain.Chat{}).Where("channel = ?", dto.Channel).Where("is_owner = ?", true).First(&chat)
	if chat == nil || chat.ID == 0 || chat.UserId != cast.ToInt(dto.CauserId) {
		return nil, errors.New("Can't Enable Chat")
	}

	err := r.db.Model(&domain.Chat{}).Where("channel = ?", dto.Channel).Update("status", consts.ACCEPT_CHAT_STATUS)
	if err.Error != nil {
		return nil, err.Error
	}

	chat.Status = consts.ACCEPT_CHAT_STATUS

	return chat, nil
}

func (r ChatRepository) RejectOffer(ctx context.Context, dto *dto.RejectOffer) (*domain.Message, error) {
	message := &domain.Message{}
	message.ID = dto.MessageId
	r.db.First(&message)
	message.TransferOfferStatus = utils.GetAsPointer(consts.REJECT_TRANSFER_OFFER_STATUS)
	r.db.Save(&message)
	return message, nil
}

func (r ChatRepository) GetChat(ctx context.Context, dto *dto.GetChat) (*domain.Chat, error) {
	var chat *domain.Chat
	r.db.Model(&domain.Chat{}).Where("channel = ?", dto.Channel).Where("user_id != ?", cast.ToInt(dto.CauserId)).First(&chat)
	if chat == nil || chat.ID == 0 {
		return nil, errors.New("Can't Enable Chat")
	}
	if dto.GetMarkAsRead() && chat.UnreadMessagesCount != 0 {
		chat.UnreadMessagesCount = 0
		r.db.Save(&chat)
	}
	return chat, nil
}

func (r ChatRepository) StoreMessage(ctx context.Context, dto *dto.SendMessage) (*domain.Message, error) {
	user, err := r.driverRepository.Find(&ctx, uint(*utils.StringToUint(dto.CauserId)))
	if err != nil {
		return nil, err
	}

	message := &domain.Message{
		Channel:                 dto.Channel,
		Message:                 dto.Message,
		MessageType:             dto.MessageType,
		BrandId:                 dto.BrandId,
		TransferId:              dto.TransferId,
		SenderId:                int64(*utils.StringToUint(dto.CauserId)),
		SenderType:              dto.CauserType,
		IsPrivate:               strings.HasPrefix(dto.Channel, "private-"),
		User:                    utils.StructToMap(user.ToResponse(), "json"),
		TransferOffersRequested: dto.TransferOffersRequested,
	}
	if dto.TransferOffersRequested {
		message.TransferOfferStatus = utils.GetAsPointer[string](consts.PENDING_TRANSFER_OFFER_STATUS)
	}

	result := r.db.Create(&message)

	go func() {
		if err := r.db.Model(&domain.Chat{}).Where("channel = ?", message.Channel).Where("user_id != ?", message.SenderId).
			UpdateColumn("unread_messages_count", gorm.Expr("unread_messages_count + ?", 1)).Error; err != nil {
			r.logger.Error(err)
		}
	}()

	go func() {
		updateResult := r.db.Model(&domain.Chat{}).Where("channel = ?", message.Channel).Updates(&domain.Chat{Status: consts.ACCEPT_CHAT_STATUS, LastMessage: map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"brand_id":     message.BrandId,
			"transfer_id":  message.TransferId,
			"message":      message.Message,
			"message_type": message.MessageType,
		}})
		if updateResult.Error != nil {
			r.logger.Error("Update Last Message Error => ", updateResult.Error.Error())
		}
	}()

	go func() {
		errIncrement := r.driverRepository.IncrementSoldTripsByValue(&ctx, uint(message.SenderId), 1)
		if errIncrement != nil {
			r.logger.Error(errIncrement)
		}
	}()

	return message, result.Error
}

func (r ChatRepository) UpdateMessage(ctx context.Context, dto *dto.UpdateMessage) (*domain.Message, error) {
	message := &domain.Message{}
	message.ID = dto.MessageId

	r.db.First(&message)

	if !r.gate.Authorize(message, gate.MethodNames.Update, ctx) {
		r.logger.Error("UpdateMessage -> UnAuthorized -> ", message.ID)
		return nil, validators.GetError(&ctx, localization.E1006, nil)
	}

	message.Message = dto.Message

	result := r.db.Save(&message)

	return message, result.Error
}

func (r ChatRepository) DeleteMessage(ctx context.Context, dto *dto.DeleteMessage) (*domain.Message, error) {
	message := &domain.Message{}
	message.ID = dto.MessageId
	r.db.First(&message)

	if !r.gate.Authorize(message, gate.MethodNames.Delete, ctx) {
		r.logger.Error("DeleteMessage -> UnAuthorized -> ", message.ID)
		return nil, validators.GetError(&ctx, localization.E1006, nil)
	}

	result := r.db.Delete(message)

	return message, result.Error
}
