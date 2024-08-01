package subscribers

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
	consts3 "samm/internal/module/notification/consts"
	notification2 "samm/internal/module/notification/dto/notification"
	"samm/internal/module/notification/gateways/onesignal/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/subscribers/notification"
	"samm/pkg/utils"
)

var SubscriberTopics = struct{ OrderChange string }{OrderChange: "order_changed"}

type OrderSubscriber struct {
	orderUseCase domain.OrderUseCase
	EventBus     EventBus.Bus
}

func OrderChangeSubscriber(eventBus EventBus.Bus, orderUseCase domain.OrderUseCase) {
	o := OrderSubscriber{orderUseCase: orderUseCase, EventBus: eventBus}

	eventBus.SubscribeAsync(SubscriberTopics.OrderChange, o.orderChange, false)
}

func (o OrderSubscriber) orderChange(order *domain.Order) {
	ctx := context.Background()
	o.pushNotificationOrder(ctx, order)
	o.orderUseCase.UpdateRealTimeDb(ctx, order)
}
func (o OrderSubscriber) pushNotificationOrder(ctx context.Context, order *domain.Order) {

	notificationData := notification2.GeneralNotification{}
	notificationData.Country = order.Location.Country.Id
	notificationData.To = []notification2.NotificationReceiver{
		{
			Id:              utils.ConvertObjectIdToStringId(order.User.ID),
			Model:           consts.UserModelType,
			LogNotification: true,
		},
		{
			Id:        utils.ConvertObjectIdToStringId(order.Location.ID),
			Model:     consts.LocationModelType,
			AccountId: utils.ConvertObjectIdToStringId(order.Location.Account.Id),
		},
	}
	notificationData.NotificationData = map[string]string{
		"location_name": order.Location.Name.En,
		"order_serial":  order.SerialNum,
	}
	notificationData.NotificationCode = notification.GetNotificationCodeBasedOnStatus(order)

	if notificationData.NotificationCode == "" {
		logrus.Print("No Suitable Notification Dode", order)

		return
	}

	o.EventBus.Publish(consts3.SEND_NOTIFICATION, notificationData)
}
