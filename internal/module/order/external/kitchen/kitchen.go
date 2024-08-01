package kitchen

import (
	kitchenDomain "samm/internal/module/kitchen/domain"
	"samm/pkg/logger"
)

type IService struct {
	KitchenUseCase kitchenDomain.KitchenUseCase
	Logger         logger.ILogger
}
