package external

import (
	menuDomain "samm/internal/module/menu/domain"
	"samm/internal/module/order/external/menu"
	"samm/internal/module/order/external/retails"
	"samm/internal/module/retails/domain"
	domain2 "samm/internal/module/user/domain"
)

type ExtService struct {
	RetailsIService retails.IService
	MenuIService    menu.IService
}

func NewExternalService(locationUseCase domain.LocationUseCase, collectionUseCase domain2.CollectionMethodUseCase, MenuUseCase menuDomain.MenuGroupItemRepository) ExtService {
	return ExtService{
		RetailsIService: retails.IService{LocationUseCase: locationUseCase, CollectionUseCase: collectionUseCase},
		MenuIService:    menu.IService{MenuUseCase: MenuUseCase},
	}
}
