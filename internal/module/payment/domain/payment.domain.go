package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/response"
	"samm/pkg/validators"
)

type Payment struct {
	mgm.DefaultModel `bson:",inline"`
	TransactionId    primitive.ObjectID `json:"transaction_id" bson:"transaction_id"`
	TransactionType  string             `json:"transaction_type" bson:"transaction_type"`
	Status           string             `json:"status" bson:"status"`
	Amount           float64            `json:"amount" bson:"amount"`
	Currency         string             `json:"currency" bson:"currency"`
	RequestId        string             `json:"request_id" bson:"request_id"`
	Request          interface{}        `json:"request" bson:"request"`
	Response         interface{}        `json:"response" bson:"response"`
	PaymentType      string             `json:"payment_type" bson:"payment_type"`
	CardType         string             `json:"card_type" bson:"card_type"`
	CardNumber       string             `json:"card_number" bson:"card_number"`
}

type PaymentUseCase interface {
	Pay(ctx context.Context, dto *payment.PayDto) (payResponse response.PayResponse, err validators.ErrorResponse)
	UpdateSession(ctx context.Context, dto *payment.UpdateSession) (payResponse response.UpdateSessionResponse, err validators.ErrorResponse)
	MyFatoorahWebhook(ctx context.Context, dto *payment.MyFatoorahWebhookPayload) (payResponse interface{}, err validators.ErrorResponse)
	AuthorizePayment(ctx context.Context, payload *payment.AuthorizePayload) (payResponse response.PayResponse, err validators.ErrorResponse)
}

type PaymentRepository interface {
	CreateTransaction(ctx context.Context, document *Payment) (response *Payment, err error)
	UpdateTransaction(ctx context.Context, document *Payment) (err error)
	FindPaymentTransaction(ctx context.Context, id string, transactionId string, transactionType string) (payment *Payment, err error)
	FindPaymentTransactionByRequestId(ctx context.Context, requestId string) (payment *Payment, err error)
}
