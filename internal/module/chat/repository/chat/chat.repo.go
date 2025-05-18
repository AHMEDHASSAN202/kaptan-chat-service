package chat

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/gate"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"math/rand"
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
	r.db.Where("user_type = ?", dto.CauserType).Where("user_id = ?", dto.CauserId).Find(&chats)
	return chats
}

func (r ChatRepository) AddPrivateChat(ctx context.Context, dto *dto.AddPrivateChat) (*domain.Chat, *domain.Message, error) {
	message := domain.Message{}
	message.ID = dto.MessageId
	r.db.First(&message)

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(999-111+1) + 111
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
	if user == nil {
		return nil, nil, errors.New("Something Went Wrong")
	}
	chat := &domain.Chat{
		Channel:             channel,
		UserType:            dto.CauserType,
		UserId:              int(*utils.StringToUint(dto.CauserId)),
		TransferId:          message.TransferId,
		UnreadMessagesCount: 1,
		IsOwner:             false,
		User:                message.User,
		LastMessage: map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"brand_id":     message.BrandId,
			"transfer_id":  message.TransferId,
			"message":      message.Message,
			"message_type": message.MessageType,
		},
		Disabled: true,
	}

	ownerUserChat := &domain.Chat{
		Channel:             channel,
		UserType:            message.SenderType,
		UserId:              int(message.SenderId),
		TransferId:          message.TransferId,
		UnreadMessagesCount: 1,
		IsOwner:             true,
		User: map[string]interface{}{
			"id":      user.ID,
			"name":    user.Name,
			"phone":   user.Phone,
			"address": user.Address,
			"image":   "",
		},
		LastMessage: map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"brand_id":     message.BrandId,
			"transfer_id":  message.TransferId,
			"message":      message.Message,
			"message_type": message.MessageType,
		},
		Disabled: true,
	}

	result := r.db.Create([]*domain.Chat{
		chat, ownerUserChat,
	})

	return chat, &message, result.Error
}

func (r ChatRepository) EnablePrivateChat(ctx context.Context, dto *dto.EnablePrivateChat) (*domain.Chat, error) {
	chat := domain.Chat{}
	r.db.Model(&domain.Chat{}).Where("channel = ?", dto.Channel).Where("is_owner = ?", true).First(&chat)
	if chat.ID == 0 || chat.UserId != cast.ToInt(dto.CauserId) {
		return nil, errors.New("Can't Enable Chat")
	}

	err := r.db.Model(&domain.Chat{}).Where("channel = ?", dto.Channel).Update("disabled", false)
	if err.Error != nil {
		return nil, err.Error
	}

	chat.Disabled = false

	return &chat, nil
}

func (r ChatRepository) StoreMessage(ctx context.Context, dto *dto.SendMessage) (*domain.Message, error) {
	message := &domain.Message{
		Channel:     dto.Channel,
		Message:     dto.Message,
		MessageType: dto.MessageType,
		BrandId:     dto.BrandId,
		TransferId:  dto.TransferId,
		SenderId:    int64(*utils.StringToUint(dto.CauserId)),
		SenderType:  dto.CauserType,
	}

	user, err := r.driverRepository.Find(&ctx, uint(*utils.StringToUint(dto.CauserId)))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Something Went Wrong")
	}

	message.User = map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"phone":   user.Phone,
		"address": user.Address,
		"image":   "",
	}
	result := r.db.Create(&message)

	go func() {
		if err := r.db.Model(&domain.Chat{}).Where("channel = ?", message.Channel).Where("user_id != ?", message.SenderId).Where("user_type != ?", message.SenderType).
			UpdateColumn("unread_messages_count", gorm.Expr("unread_messages_count + ?", 1)).Error; err != nil {
			r.logger.Error(err)
		}
	}()

	go func() {
		updateResult := r.db.Model(&domain.Chat{}).Where("channel = ?", message.Channel).Updates(&domain.Chat{LastMessage: map[string]interface{}{
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
