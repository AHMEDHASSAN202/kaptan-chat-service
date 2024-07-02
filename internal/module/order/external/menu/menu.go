package menu

import (
	menuDomain "samm/internal/module/menu/domain"
)

type IService struct {
	MenuUseCase menuDomain.MenuGroupItemRepository
}
