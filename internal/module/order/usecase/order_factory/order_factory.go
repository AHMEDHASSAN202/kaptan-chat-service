package order_factory

import (
	"errors"
	"github.com/asaskevich/EventBus"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/external"
	"samm/pkg/database/redis"
	"samm/pkg/gate"
	"samm/pkg/logger"
)

type OrderFactory struct {
	orderTypes map[string]func() IOrder
}

func NewOrderFactory(validator *validator.Validate, eventBus EventBus.Bus, extService external.ExtService, logger logger.ILogger, orderRepo domain.OrderRepository, redisClient *redis.RedisClient, gate *gate.Gate) *OrderFactory {
	return &OrderFactory{
		orderTypes: map[string]func() IOrder{
			"ktha": func() IOrder {
				return &KthaOrder{Deps: Deps{eventBus: eventBus, validator: validator, extService: extService, logger: logger, orderRepo: orderRepo, redisClient: redisClient, gate: gate}}
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
