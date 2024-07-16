package order

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderPaidDto struct {
	OrderId       primitive.ObjectID
	TransactionId string `header:"transaction_id" bson:"transaction_id"`
	PaymentType   string `json:"payment_type" bson:"payment_type"`
	CardType      string `json:"card_type" bson:"card_type"`
	CardNumber    string `json:"card_number" bson:"card_number"`
}
