package domain

import (
	"context"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/gateways/myfatoorah/requests"
	"samm/internal/module/payment/gateways/myfatoorah/responses"
	"samm/pkg/validators"
)

type MyFatoorahService interface {
	PayCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *Payment) (paymentResponse responses.ExecutePaymentCardResponse, requestPayload requests.ExecutePaymentCardRequest, invoiceId int, err validators.ErrorResponse)
	ExecutePaymentCard(ctx context.Context, paymentTransaction *Payment) (paymentResponse responses.ExecutePaymentCardResponse, requestPayload requests.ExecutePaymentCardRequest, invoiceId int, err validators.ErrorResponse)
	InitPaymentCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *Payment) (paymentResponse responses.InitSessionResponse, redirectUrl string, err validators.ErrorResponse)
	FindPayment(ctx context.Context, invoiceId string) (paymentResponse responses.GetPaymentStatusResponse, err validators.ErrorResponse)
	UpdatePaymentStatus(ctx context.Context, invoiceId string, capture bool) (err validators.ErrorResponse)
	ApplePay(ctx context.Context, dto *payment.PayDto, paymentTransaction *Payment) (paymentResponse responses.ExecutePaymentResponse, requestPayload requests.ApplePayExecutePaymentCardRequest, err validators.ErrorResponse)
}
