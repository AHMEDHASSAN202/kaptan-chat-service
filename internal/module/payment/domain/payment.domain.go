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
	TransactionId    primitive.ObjectID `json:"transaction_id"`
	TransactionType  string             `json:"transaction_type"`
	Status           string             `json:"status"`
	Amount           float64            `json:"amount"`
	Currency         string             `json:"currency"`
	Request          interface{}        `json:"request"`
	Response         interface{}        `json:"response"`
	PaymentType      string             `json:"payment_type"`
	CardType         string             `json:"card_type"`
}

type PaymentUseCase interface {
	Pay(ctx context.Context, dto *payment.PayDto) (payResponse response.PayResponse, err validators.ErrorResponse)
}

type PaymentRepository interface {
	CreateTransaction(ctx context.Context, document *Payment) (response *Payment, err error)
	UpdateTransaction(ctx context.Context, document *Payment) (err error)
	FindPaymentTransaction(ctx context.Context, id string, transactionId string, transactionType string) (payment *Payment, err error)
}
