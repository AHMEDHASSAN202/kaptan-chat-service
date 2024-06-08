package location

import (
	"samm/internal/module/retails/domain"
	"samm/pkg/logger"
)

type LocationUseCase struct {
	repo   domain.LocationRepository
	logger logger.ILogger
}

const tag = " LocationUseCase "

func NewLocationUseCase(repo domain.LocationRepository, logger logger.ILogger) domain.LocationUseCase {
	return &LocationUseCase{
		repo:   repo,
		logger: logger,
	}
}
