package menu

import (
	menuDomain "samm/internal/module/menu/domain"
	"samm/pkg/logger"
)

type IService struct {
	MenuUseCase menuDomain.MenuGroupItemRepository
	Logger      logger.ILogger
}
