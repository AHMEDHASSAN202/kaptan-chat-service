package order_factory

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/order/external"
	"samm/pkg/logger"
)

type OrderFactory struct {
	orderTypes map[string]func() IOrder
}

func NewOrderFactory(validator *validator.Validate, extService *external.ExtService, logger logger.ILogger) *OrderFactory {
	return &OrderFactory{
		orderTypes: map[string]func() IOrder{
			"ktha": func() IOrder {
				return &NgoOrder{Deps: Deps{validator: validator, extService: extService, logger: logger}}
			},
		},
	}
}

func (f *OrderFactory) Make(orderType string) (IOrder, error) {
	if orderFunc, exists := f.orderTypes[orderType]; exists {
		return orderFunc().Make(), nil
	}
	return nil, errors.New("unknown order type")
}
