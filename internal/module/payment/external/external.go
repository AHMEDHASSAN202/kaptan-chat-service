package external

import (
	"samm/internal/module/order/domain"
	"samm/internal/module/payment/external/order"
)

type ExtService struct {
	OrderService order.IService
}

func NewExternalService() *ExtService {
	return &ExtService{}
}
func (e *ExtService) SetOrderUseCase(orderUseCase domain.OrderUseCase) {
	e.OrderService = order.IService{OrderUseCase: orderUseCase}
}
