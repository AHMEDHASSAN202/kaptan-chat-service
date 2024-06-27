package domain

import (
	"context"
	"samm/internal/module/payment/dto/payment"
	"samm/pkg/validators"
)

type MyFatoorahService interface {
	PayCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *Payment) (err validators.ErrorResponse)
}
