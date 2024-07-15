package payment

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/order/external/payment/responses"
	"samm/internal/module/payment/dto/payment"
	"samm/pkg/validators"
)

func (i IService) AuthorizePayment(ctx context.Context, transactionId string, capture bool) (responses.PaymentResponse, validators.ErrorResponse) {

	input := payment.AuthorizePayload{
		Capture:       capture,
		TransactionId: transactionId,
	}
	var paymentResponse responses.PaymentResponse
	paymentDomain, err := i.PaymentUseCase.AuthorizePayment(ctx, &input)
	if err.IsError {
		return responses.PaymentResponse{}, err
	}

	copier.Copy(&paymentResponse, &paymentDomain)

	return paymentResponse, validators.ErrorResponse{}
}
