package payment

import (
	"samm/internal/module/payment/domain"
)

type IService struct {
	PaymentUseCase domain.PaymentUseCase
}
