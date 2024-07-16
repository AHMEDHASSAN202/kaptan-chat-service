package external

import (
	"samm/internal/module/order/domain"
	"samm/internal/module/payment/external/order"
)

type ExtService struct {
	OrderService order.IService
}

func NewExternalService(orderUseCase domain.OrderUseCase) ExtService {
	return ExtService{
		OrderService: order.IService{OrderUseCase: orderUseCase},
	}
}
