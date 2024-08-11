package payment

import (
	"context"
	"fmt"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/external"
	"samm/internal/module/payment/response"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strconv"
)

type PaymentUseCase struct {
	repo              domain.PaymentRepository
	myfatoorahService domain.MyFatoorahService
	cardRepo          domain.CardRepository
	logger            logger.ILogger
	extService        *external.ExtService
}

func (p PaymentUseCase) GetPaymentStatus(ctx context.Context, dto *payment.GetPaymentStatus) (payResponse *domain.Payment, err validators.ErrorResponse) {
	paymentTransaction, errResp := p.repo.FindPaymentTransaction(ctx, "", dto.TransactionId, dto.TransactionType)
	if errResp != nil {
		return nil, validators.GetErrorResponseFromErr(errResp)
	}

	if paymentTransaction.Status == consts.PaymentPaidStatus || paymentTransaction.Status == consts.PaymentHoldStatus {
		return paymentTransaction, validators.ErrorResponse{}
	}
	fmt.Println(paymentTransaction.Status)

	// Call Mf To check the status
	responsePay, errRe := p.myfatoorahService.FindPayment(ctx, paymentTransaction.RequestId)

	// Update Payment Transaction to Paid or failed
	if !errRe.IsError && responsePay.Data.InvoiceStatus == consts.MyFatoorahPaidInvoiceStatus {

		paymentTransaction.Status = consts.PaymentPaidStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
			return nil, validators.GetErrorResponseFromErr(errUpdate)
		}
		UpdateOrderStatus(p, ctx, paymentTransaction)

		return paymentTransaction, validators.ErrorResponse{}
	}

	// Check The Hold Status
	hold := len(responsePay.Data.InvoiceTransactions) > 0 && responsePay.Data.InvoiceTransactions[0].TransactionStatus == consts.MyFatoorahHoldInvoiceStatus
	if !errRe.IsError && hold {

		paymentTransaction.Status = consts.PaymentHoldStatus
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway

		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
			return nil, validators.GetErrorResponseFromErr(errUpdate)

		}
		UpdateOrderStatus(p, ctx, paymentTransaction)

		return paymentTransaction, validators.ErrorResponse{}
	}

	paymentTransaction.Status = consts.PaymentFailedStatus
	paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
	if len(responsePay.Data.InvoiceTransactions) > 0 {
		paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway
	}

	errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
	if errUpdate != nil {
		p.logger.Error("Update Error => ", errUpdate)
		return nil, validators.GetErrorResponseFromErr(errUpdate)

	}
	transactionError := localization.PaymentError

	if len(responsePay.Data.InvoiceTransactions) > 0 && responsePay.Data.InvoiceTransactions[0].ErrorCode != "" {

		transactionError = responsePay.Data.InvoiceTransactions[0].ErrorCode
	}
	return nil, validators.GetErrorResponseWithErrors(&ctx, transactionError, nil)

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
func (p PaymentUseCase) UpdateSession(ctx context.Context, dto *payment.UpdateSession) (payResponse response.UpdateSessionResponse, err validators.ErrorResponse) {
	// Find Transaction
	paymentTransaction, errResp := p.repo.FindPaymentTransactionByRequestId(ctx, dto.SessionId)
	if errResp != nil {
		return payResponse, validators.GetErrorResponseFromErr(errResp)
	}
	executeResponse, executeRequest, innvoiceId, errRes := p.myfatoorahService.ExecutePaymentCard(ctx, paymentTransaction)
	if errRes.IsError {
		paymentTransaction.Status = consts.PaymentFailedStatus
		paymentTransaction.Request = utils.ConvertStructToMap(executeRequest)
		paymentTransaction.Response = utils.ConvertStructToMap(executeResponse)
		paymentTransaction.RequestId = strconv.Itoa(innvoiceId)
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
		}
		return payResponse, errRes
	}
	paymentTransaction.Request = utils.ConvertStructToMap(executeRequest)
	paymentTransaction.Response = utils.ConvertStructToMap(executeResponse)
	paymentTransaction.RequestId = strconv.Itoa(innvoiceId)
	errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
	if errUpdate != nil {
		p.logger.Error("Update Error => ", errUpdate)
	}
	payResponse.RedirectUrl = &executeResponse.Data.PaymentURL

	// handle update user cards for user
	HandleUpdateUserCards(p, ctx, paymentTransaction)

	return payResponse, err
}
