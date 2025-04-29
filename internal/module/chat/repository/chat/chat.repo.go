package chat

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/pkg/gate"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
)

type ChatRepository struct {
	db     *gorm.DB
	logger logger.ILogger
	gate   *gate.Gate
}

func NewChatRepository(logger logger.ILogger, db *gorm.DB, gate *gate.Gate) domain.ChatRepository {
	return &ChatRepository{
		db:     db,
		logger: logger,
		gate:   gate,
	}
}

func (r ChatRepository) StoreMessage(ctx context.Context, dto *dto.SendMessage) (*domain.Message, error) {
	message := &domain.Message{
		Channel:     dto.Channel,
		Message:     dto.Message,
		MessageType: dto.MessageType,
		BrandId:     sql.NullInt64{Int64: *dto.BrandId},
		SenderId:    int64(*utils.StringToUint(dto.CauserId)),
		SenderType:  dto.CauserType,
	}

	if dto.TransferId != nil {
		message.TransferId = sql.NullInt64{Int64: *dto.TransferId}
	}

	if dto.OwnerTransferId != nil {
		message.OwnerTransferId = sql.NullInt64{Int64: *dto.OwnerTransferId}
	}

	result := r.db.Create(&message)

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
