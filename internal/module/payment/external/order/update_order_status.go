package order

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/external/order/responses"
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
func (i IService) FindOrder(ctx context.Context, orderId string) (orderResponse responses.OrderResponse, err validators.ErrorResponse) {

	orderDomain, err := i.OrderUseCase.FindOrderForDashboard(&ctx, orderId)
	if err.IsError {
		return orderResponse, err
	}
	copier.Copy(&orderResponse, orderDomain)
	return orderResponse, err
}
