package payment

import (
	"context"
	"samm/internal/module/order/external/payment/responses"
	"samm/pkg/validators"
)

type Interface interface {
	AuthorizePayment(ctx context.Context, transactionId string, capture bool) (responses.PaymentResponse, validators.ErrorResponse)
}
