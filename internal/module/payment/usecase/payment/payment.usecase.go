package payment

import (
	"context"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/response"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type PaymentUseCase struct {
	repo              domain.PaymentRepository
	myfatoorahService domain.MyFatoorahService
	logger            logger.ILogger
}

func (p PaymentUseCase) Pay(ctx context.Context, dto *payment.PayDto) (paymentResponse response.PayResponse, err validators.ErrorResponse) {
	// Check The Type and pay with
	switch dto.PaymentType {
	case consts.PaymentTypeCard:
		return PayCard(p, ctx, dto)
	default:
		return PayApplePay(p, ctx, dto)
	}
	return paymentResponse, err
}

func NewPaymentUseCase(repo domain.PaymentRepository, myfatoorahService domain.MyFatoorahService, logger logger.ILogger) domain.PaymentUseCase {
	return &PaymentUseCase{
		repo:              repo,
		myfatoorahService: myfatoorahService,

		logger: logger,
	}
}
