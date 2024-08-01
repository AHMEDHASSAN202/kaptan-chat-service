package notification

import (
	consts3 "samm/internal/module/notification/consts"
	consts2 "samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
)

func GetNotificationCodeBasedOnStatus(order *domain.Order) (notificationCode string) {
	switch order.Status {
	case consts2.OrderStatus.Pending:
		return consts3.ORDER_PENDING
	case consts2.OrderStatus.TimedOut:
		return consts3.ORDER_TIMEOUT
	case consts2.OrderStatus.Accepted:
		if order.ArrivedAt != nil {
			return consts3.ORDER_ARRIVED
		}
		return consts3.ORDER_ACCEPTED
	case consts2.OrderStatus.Cancelled:
		return consts3.ORDER_CANCELLED
	case consts2.OrderStatus.Rejected:
		return consts3.ORDER_REJECTED
	case consts2.OrderStatus.Rejected:
		return consts3.ORDER_REJECTED
	case consts2.OrderStatus.ReadyForPickup:
		if order.ArrivedAt != nil {
			return consts3.ORDER_ARRIVED
		}
		return consts3.ORDER_READY_PICKUP
	case consts2.OrderStatus.PickedUp:
		return consts3.ORDER_PICKUP
	case consts2.OrderStatus.NoShow:
		if order.ArrivedAt != nil {
			return consts3.ORDER_ARRIVED
		}
		return consts3.ORDER_NOSHOW
	}
	return
}
