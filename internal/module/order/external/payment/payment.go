package payment

import (
	"samm/internal/module/payment/domain"
	"samm/pkg/logger"
)

type IService struct {
	PaymentUseCase domain.PaymentUseCase
	Logger         logger.ILogger
}
