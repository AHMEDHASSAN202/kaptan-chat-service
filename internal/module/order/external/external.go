package external

import (
	menuDomain "samm/internal/module/menu/domain"
	"samm/internal/module/order/external/menu"
	"samm/internal/module/order/external/payment"
	"samm/internal/module/order/external/retails"
	domain3 "samm/internal/module/payment/domain"
	"samm/internal/module/retails/domain"
	domain2 "samm/internal/module/user/domain"
	"samm/pkg/logger"
)

type ExtService struct {
	logger          logger.ILogger
	RetailsIService retails.IService
	MenuIService    menu.IService
	PaymentIService payment.IService
}

func NewExternalService(locationUseCase domain.LocationUseCase, logger logger.ILogger, collectionUseCase domain2.CollectionMethodUseCase, MenuUseCase menuDomain.MenuGroupItemRepository, paymentUseCase domain3.PaymentUseCase, accountUseCase domain.AccountUseCase) ExtService {
	return ExtService{
		RetailsIService: retails.IService{LocationUseCase: locationUseCase, CollectionUseCase: collectionUseCase, AccountUseCase: accountUseCase, Logger: logger},
		MenuIService:    menu.IService{MenuUseCase: MenuUseCase, Logger: logger},
		PaymentIService: payment.IService{PaymentUseCase: paymentUseCase, Logger: logger},
		logger:          logger,
	}
}
