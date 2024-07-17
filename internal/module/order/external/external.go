package external

import (
	menuDomain "samm/internal/module/menu/domain"
	"samm/internal/module/order/external/menu"
	"samm/internal/module/order/external/payment"
	"samm/internal/module/order/external/retails"
	domain3 "samm/internal/module/payment/domain"
	"samm/internal/module/retails/domain"
	domain2 "samm/internal/module/user/domain"
)

type ExtService struct {
	RetailsIService retails.IService
	MenuIService    menu.IService
	PaymentIService payment.IService
}

func NewExternalService(locationUseCase domain.LocationUseCase, collectionUseCase domain2.CollectionMethodUseCase, MenuUseCase menuDomain.MenuGroupItemRepository, paymentUseCase domain3.PaymentUseCase, accountUseCase domain.AccountUseCase) ExtService {
	return ExtService{
		RetailsIService: retails.IService{LocationUseCase: locationUseCase, CollectionUseCase: collectionUseCase, AccountUseCase: accountUseCase},
		MenuIService:    menu.IService{MenuUseCase: MenuUseCase},
		PaymentIService: payment.IService{PaymentUseCase: paymentUseCase},
	}
}
