package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentResponse struct {
	Id              primitive.ObjectID `json:"id"`
	TransactionId   primitive.ObjectID `json:"transaction_id"`
	TransactionType string             `json:"transaction_type"`
	Status          string             `json:"status"`
	Amount          float64            `json:"amount"`
	Currency        string             `json:"currency"`
	RequestId       string             `json:"request_id"`
	Request         interface{}        `json:"request"`
	Response        interface{}        `json:"response"`
	PaymentType     string             `json:"payment_type"`
	CardType        string             `json:"card_type"`
}
