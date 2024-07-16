package order

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/payment/domain"
	"samm/pkg/validators"
)

type Interface interface {
	SetOrderPaid(ctx context.Context, orderId primitive.ObjectID, payment domain.Payment) validators.ErrorResponse
}
