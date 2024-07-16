package order

import "samm/internal/module/order/domain"

type IService struct {
	OrderUseCase domain.OrderUseCase
}
