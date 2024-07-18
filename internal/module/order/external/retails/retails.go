package retails

import (
	"samm/internal/module/retails/domain"
	domain2 "samm/internal/module/user/domain"
	"samm/pkg/logger"
)

type IService struct {
	LocationUseCase   domain.LocationUseCase
	AccountUseCase    domain.AccountUseCase
	CollectionUseCase domain2.CollectionMethodUseCase
	Logger            logger.ILogger
}
