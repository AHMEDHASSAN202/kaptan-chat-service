package payment

import (
	"context"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/response"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

func CreateTransactionBuilder(dto *payment.PayDto) *domain.Payment {
	paymentDomain := domain.Payment{}
	paymentDomain.TransactionId = utils.ConvertStringIdToObjectId(dto.TransactionId)
	paymentDomain.TransactionType = dto.TransactionType
	paymentDomain.Amount = 20      // Todo Depend On Find Order
	paymentDomain.Currency = "SAR" // Todo Depend On Find Order
	paymentDomain.PaymentType = dto.PaymentType
	paymentDomain.Status = consts.PaymentPendingStatus

	return &paymentDomain
}

func PayCard(p PaymentUseCase, ctx context.Context, dto *payment.PayDto) (response response.PayResponse, err validators.ErrorResponse) {
	paymentTransaction, errResp := p.repo.FindPaymentTransaction(ctx, "", dto.TransactionId, dto.TransactionType)
	if errResp == nil && (paymentTransaction.Status == consts.PaymentPaidStatus || paymentTransaction.Status == consts.PaymentHoldStatus) {
		response.Transaction = paymentTransaction

		return response, validators.ErrorResponse{}
	}
	// Check if it is pending then call mf to get payment status
	//if errResp == nil && (paymentTransaction.Status == consts.PaymentPendingStatus) {
	//	// Call Mf To check the status
	//	response.Transaction = paymentTransaction
	//	return response, validators.ErrorResponse{}
	//}

	// Create The Transaction
	paymentTransaction, errResp = p.repo.CreateTransaction(ctx, CreateTransactionBuilder(dto))
	if errResp != nil {

		return response, validators.GetErrorResponseFromErr(errResp)
	}

	// Call Myfatoorah
	err = p.myfatoorahService.PayCard(ctx, dto, paymentTransaction)
	// Check if it using token
	//
	return response, err
}
func PayApplePay(p PaymentUseCase, ctx context.Context, dto *payment.PayDto) (paymentResponse response.PayResponse, err validators.ErrorResponse) {
	return paymentResponse, err
}
