package subscribers

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
	"samm/internal/module/order/domain"
)

var SubscriberTopics = struct{ OrderChange string }{OrderChange: "order_changed"}

type OrderSubscriber struct {
	orderUseCase domain.OrderUseCase
}

func OrderChangeSubscriber(eventBus EventBus.Bus, orderUseCase domain.OrderUseCase) {
	o := OrderSubscriber{orderUseCase: orderUseCase}

	eventBus.Subscribe(SubscriberTopics.OrderChange, o.orderChange)
}

func (o OrderSubscriber) orderChange(order *domain.Order) {
	ctx := context.Background()
	logrus.Print("an event has triggered by an order change", order)

	o.orderUseCase.UpdateRealTimeDb(ctx, order)
}
