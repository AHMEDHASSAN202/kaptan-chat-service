package repository

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/transfer/domain"
	"kaptan/pkg/logger"
)

type Repository struct {
	logger logger.ILogger
	db     *gorm.DB
}

func NewTransferRepository(log logger.ILogger, db *gorm.DB) domain.TransferRepository {
	return &Repository{
		logger: log,
		db:     db,
	}
}

func (r *Repository) Find(ctx *context.Context, id uint) (*domain.Transfer, error) {
	transfer := domain.Transfer{ID: id}
	r.db.First(&transfer)
	return &transfer, nil
}
