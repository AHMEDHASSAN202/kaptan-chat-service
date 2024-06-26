package payment

import (
	"samm/internal/module/payment/domain"
	"samm/pkg/logger"
)

type PaymentUseCase struct {
	repo   domain.PaymentRepository
	logger logger.ILogger
}

func NewPaymentUseCase(repo domain.PaymentRepository, logger logger.ILogger) domain.PaymentUseCase {
	return &PaymentUseCase{
		repo:   repo,
		logger: logger,
	}
}
