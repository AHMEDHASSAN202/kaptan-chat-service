package repository

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/transfer/domain"
	"kaptan/internal/module/transfer/dto"
	"kaptan/pkg/logger"
	"time"
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
	result := r.db.First(&transfer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transfer, nil
}

func (r *Repository) AssignSellerToTransfer(ctx *context.Context, driverID uint, transferID uint) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferID)
	if err != nil {
		return nil, err
	}

	// Step 2: Assign driverID as sellerID
	transfer.SellerID = &driverID // assuming SellerID is a *uint in the domain

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}

func (r *Repository) MarkTransferAsStart(ctx *context.Context, transferDto *dto.StartTransfer) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferDto.TransferId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Step 2: Change status to "Start"
	transfer.Status = domain.TransferStatusStart // assuming Status is a string
	transfer.StartAt = &now

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}

func (r *Repository) MarkTransferAsEnd(ctx *context.Context, transferDto *dto.EndTransfer) (*domain.Transfer, error) {
	// Step 1: Find the transfer
	transfer, err := r.Find(ctx, transferDto.TransferId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// Step 2: Change status to "End"
	transfer.Status = domain.TransferStatusEnd // assuming Status is a string
	transfer.EndAt = &now

	// Step 3: Save the updated transfer
	if err := r.db.Save(&transfer).Error; err != nil {
		return nil, err
	}

	// Step 4: Return updated transfer
	return transfer, nil
}
