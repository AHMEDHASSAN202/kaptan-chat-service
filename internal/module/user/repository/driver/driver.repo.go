package driver

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/user/domain"
	"kaptan/pkg/logger"
)

type Repository struct {
	logger logger.ILogger
	db     *gorm.DB
}

func NewDriverRepository(log logger.ILogger, db *gorm.DB) domain.DriverRepository {
	return &Repository{
		logger: log,
		db:     db,
	}
}

func (r *Repository) Find(ctx *context.Context, id uint) (*domain.Driver, error) {
	driver := domain.Driver{ID: id}
	r.db.First(&driver)
	return &driver, nil
}
