package driver

import (
	"context"
	"kaptan/internal/module/user/domain"
	"kaptan/pkg/logger"
)

type Repository struct {
	logger logger.ILogger
}

func NewDriverRepository(log logger.ILogger) domain.DriverRepository {
	return &Repository{
		logger: log,
	}
}

func (r *Repository) FindByToken(ctx *context.Context, token string) (*domain.User, error) {
	return nil, nil
}
