package order

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/payment/domain"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

func (i IService) SetOrderPaid(ctx context.Context, orderId primitive.ObjectID, payment domain.Payment) validators.ErrorResponse {

	return i.OrderUseCase.SetOrderPaid(ctx, &order.OrderPaidDto{
		OrderId:       orderId,
		TransactionId: utils.ConvertObjectIdToStringId(payment.ID),
		PaymentType:   payment.PaymentType,
		CardNumber:    payment.CardNumber,
		CardType:      payment.CardType,
	})
}
