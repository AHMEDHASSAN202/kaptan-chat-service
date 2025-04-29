package user

import (
	"kaptan/internal/module/user/domain"
	"kaptan/pkg/logger"
)

type Middlewares struct {
	driverRepository domain.DriverRepository
	logger           logger.ILogger
}

func NewUserMiddlewares(driverRepository domain.DriverRepository, logger logger.ILogger) *Middlewares {
	return &Middlewares{
		driverRepository: driverRepository,
		logger:           logger,
	}
}
