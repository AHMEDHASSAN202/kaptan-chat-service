package payment

import (
	"context"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/external"
	"samm/internal/module/payment/response"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"strconv"
)

type PaymentUseCase struct {
	repo              domain.PaymentRepository
	myfatoorahService domain.MyFatoorahService
	cardRepo          domain.CardRepository
	logger            logger.ILogger
	extService        *external.ExtService
}

func NewPaymentUseCase(repo domain.PaymentRepository, cardRepo domain.CardRepository, myfatoorahService domain.MyFatoorahService, logger logger.ILogger, extService *external.ExtService) domain.PaymentUseCase {
	return &PaymentUseCase{
		repo:              repo,
		myfatoorahService: myfatoorahService,
		cardRepo:          cardRepo,
		logger:            logger,
		extService:        extService,
	}
}

func (p PaymentUseCase) AuthorizePayment(ctx context.Context, payload *payment.AuthorizePayload) (payResponse response.PayResponse, err validators.ErrorResponse) {

	paymentTransaction, errResp := p.repo.FindPaymentTransaction(ctx, payload.TransactionId, "", "")
	if errResp != nil {
		return payResponse, validators.GetErrorResponseFromErr(errResp)
	}
	// Check The Status
	if paymentTransaction.Status != consts.PaymentHoldStatus {
		payResponse.Transaction = paymentTransaction
		return payResponse, err
	}

	err = p.myfatoorahService.UpdatePaymentStatus(ctx, paymentTransaction.RequestId, payload.Capture)
	if err.IsError {
		return payResponse, err
	}

	responsePay, errRe := p.myfatoorahService.FindPayment(ctx, paymentTransaction.RequestId)
	if errRe.IsError {
		return payResponse, errRe
	}
	if payload.Capture {
		responsePay.Data.InvoiceStatus = consts.MyFatoorahPaidInvoiceStatus

		paymentTransaction.Status = consts.PaymentPaidStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway

	} else {
		responsePay.Data.InvoiceStatus = consts.MyFatoorahHoldInvoiceStatus

		paymentTransaction.Status = consts.PaymentFailedStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
	}

	_ = p.repo.UpdateTransaction(ctx, paymentTransaction)
	payResponse.Transaction = paymentTransaction

	return payResponse, err
}

func (p PaymentUseCase) MyFatoorahWebhook(ctx context.Context, dto *payment.MyFatoorahWebhookPayload) (payResponse interface{}, err validators.ErrorResponse) {

	responsePay, errRe := p.myfatoorahService.FindPayment(ctx, strconv.Itoa(dto.Data.InvoiceId))
	if errRe.IsError {
		return payResponse, errRe
	}

	paymentTransaction, errResp := p.repo.FindPaymentTransaction(ctx, responsePay.Data.UserDefinedField, "", "")
	if errResp != nil {
		return payResponse, errRe
	}

	if paymentTransaction.Status != consts.PaymentPendingStatus {
		return nil, validators.ErrorResponse{}
	}

	hold := responsePay.Data.InvoiceTransactions[0].TransactionStatus == consts.MyFatoorahHoldInvoiceStatus

	if responsePay.Data.InvoiceStatus == consts.MyFatoorahPaidInvoiceStatus {

		paymentTransaction.Status = consts.PaymentPaidStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
			return payResponse, validators.GetErrorResponseFromErr(errUpdate)
		}

	} else if hold {

		paymentTransaction.Status = consts.PaymentHoldStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
			return payResponse, validators.GetErrorResponseFromErr(errUpdate)
		}

	} else {
		paymentTransaction.Status = consts.PaymentFailedStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
			return payResponse, validators.GetErrorResponseFromErr(errUpdate)
		}
	}

	paymentTransaction, errResp = p.repo.FindPaymentTransaction(ctx, responsePay.Data.UserDefinedField, "", "")
	if errResp != nil {
		return nil, validators.GetErrorResponseFromErr(errResp)
	}
	UpdateOrderStatus(p, ctx, paymentTransaction)
	return payResponse, validators.ErrorResponse{}
}

func (p PaymentUseCase) Pay(ctx context.Context, dto *payment.PayDto) (paymentResponse response.PayResponse, err validators.ErrorResponse) {
	// Check The Type and pay with
	switch dto.PaymentType {
	case consts.PaymentTypeCard:
		return PayCard(p, ctx, dto)
	default:
		return PayApplePay(p, ctx, dto)
	}
}
