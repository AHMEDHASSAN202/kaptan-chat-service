package external

import (
	menuDomain "samm/internal/module/menu/domain"
	"samm/internal/module/order/external/menu"
	"samm/internal/module/order/external/retails"
	"samm/internal/module/retails/domain"
)

type ExtService struct {
	RetailsIService retails.IService
	MenuIService    menu.IService
}

func NewExternalService(locationUseCase domain.LocationUseCase, MenuUseCase menuDomain.MenuGroupItemRepository) ExtService {
	return ExtService{
		RetailsIService: retails.IService{LocationUseCase: locationUseCase},
		MenuIService:    menu.IService{MenuUseCase: MenuUseCase},
	}
}
