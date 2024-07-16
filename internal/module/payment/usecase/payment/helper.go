package payment

import (
	"context"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/response"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"strconv"
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
	if errResp == nil && (paymentTransaction.Status == consts.PaymentPendingStatus) {
		// Call Mf To check the status
		responsePay, errRe := p.myfatoorahService.FindPayment(ctx, paymentTransaction.RequestId)

		// Update Payment Transaction to Paid or failed
		if !errRe.IsError && responsePay.Data.InvoiceStatus == consts.MyFatoorahPaidInvoiceStatus {

			paymentTransaction.Status = consts.PaymentPaidStatus
			errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
			if errUpdate != nil {
				p.logger.Error("Update Error => ", errUpdate)
				return response, validators.GetErrorResponseFromErr(errUpdate)

			}
			response.Transaction = paymentTransaction
			return response, validators.ErrorResponse{}
		}

		// Check The Hold Status
		hold := responsePay.Data.InvoiceTransactions[0].TransactionStatus == consts.MyFatoorahHoldInvoiceStatus
		if !errRe.IsError && hold {

			paymentTransaction.Status = consts.PaymentHoldStatus
			paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
			paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway

			errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
			if errUpdate != nil {
				p.logger.Error("Update Error => ", errUpdate)
				return response, validators.GetErrorResponseFromErr(errUpdate)

			}
			UpdateOrderStatus(p, ctx, paymentTransaction)

			response.Transaction = paymentTransaction
			return response, validators.ErrorResponse{}
		}
	}

	// Create The Transaction
	paymentTransaction, errResp = p.repo.CreateTransaction(ctx, CreateTransactionBuilder(dto))
	if errResp != nil {
		return response, validators.GetErrorResponseFromErr(errResp)
	}

	// Call Myfatoorah
	payResponse, payRequest, invoiceId, err := p.myfatoorahService.PayCard(ctx, dto, paymentTransaction)

	if err.IsError {
		paymentTransaction.Status = consts.PaymentFailedStatus
		paymentTransaction.Request = utils.ConvertStructToMap(payRequest)
		paymentTransaction.Response = utils.ConvertStructToMap(payResponse)
		paymentTransaction.RequestId = strconv.Itoa(invoiceId)
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
		}
		return response, err
	}
	// Check if it using token
	if dto.SaveCard && dto.PaymentToken != "" {
		// Need To Store Card Details with token
		var cardDomain domain.Card
		cardDomain.MFToken = payResponse.Data.Token
		cardDomain.Number = utils.MaskCard(dto.Card.Number)
		cardDomain.Type = dto.Card.Type
		cardDomain.UserId = utils.ConvertStringIdToObjectId(dto.UserId) // User ID
		errRe := p.cardRepo.StoreCard(ctx, &cardDomain)
		if errRe != nil {
			p.logger.Error("Save Card Error => ", errRe)
		}
	}

	paymentTransaction.Request = utils.ConvertStructToMap(payRequest)
	paymentTransaction.Response = utils.ConvertStructToMap(payResponse)
	paymentTransaction.RequestId = strconv.Itoa(invoiceId)
	errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
	if errUpdate != nil {
		p.logger.Error("Update Error => ", errUpdate)
	}

	response.RedirectUrl = &payResponse.Data.PaymentURL
	response.Transaction = paymentTransaction
	return response, err
}

func PayApplePay(p PaymentUseCase, ctx context.Context, dto *payment.PayDto) (paymentResponse response.PayResponse, err validators.ErrorResponse) {

	// Check Old Order Transaction Status
	paymentTransaction, errResp := p.repo.FindPaymentTransaction(ctx, "", dto.TransactionId, dto.TransactionType)
	if errResp == nil && (paymentTransaction.Status == consts.PaymentPaidStatus || paymentTransaction.Status == consts.PaymentHoldStatus) {
		paymentResponse.Transaction = paymentTransaction
		return paymentResponse, validators.ErrorResponse{}
	}

	// Check if it is pending then call mf to get payment status
	if errResp == nil && (paymentTransaction.Status == consts.PaymentPendingStatus) {
		// Call Mf To check the status
		responsePay, errRe := p.myfatoorahService.FindPayment(ctx, paymentTransaction.RequestId)

		// Update Payment Transaction to Paid or failed
		if !errRe.IsError && responsePay.Data.InvoiceStatus == consts.MyFatoorahPaidInvoiceStatus {

			paymentTransaction.Status = consts.PaymentPaidStatus
			errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
			if errUpdate != nil {
				p.logger.Error("Update Error => ", errUpdate)
				return paymentResponse, validators.GetErrorResponseFromErr(errUpdate)

			}
			paymentResponse.Transaction = paymentTransaction
			return paymentResponse, validators.ErrorResponse{}
		}

		// Check The Hold Status
		hold := responsePay.Data.InvoiceTransactions[0].TransactionStatus == consts.MyFatoorahHoldInvoiceStatus
		if !errRe.IsError && hold {

			paymentTransaction.Status = consts.PaymentHoldStatus
			paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
			paymentTransaction.CardType = responsePay.Data.InvoiceTransactions[0].PaymentGateway

			errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
			if errUpdate != nil {
				p.logger.Error("Update Error => ", errUpdate)
				return paymentResponse, validators.GetErrorResponseFromErr(errUpdate)

			}
			UpdateOrderStatus(p, ctx, paymentTransaction)

			paymentResponse.Transaction = paymentTransaction
			return paymentResponse, validators.ErrorResponse{}
		}
	}

	// Create The Transaction
	paymentTransaction, errResp = p.repo.CreateTransaction(ctx, CreateTransactionBuilder(dto))
	if errResp != nil {
		return paymentResponse, validators.GetErrorResponseFromErr(errResp)
	}

	responsePay, request, errRe := p.myfatoorahService.ApplePay(ctx, dto, paymentTransaction)

	if !errRe.IsError {

		paymentTransaction.Request = utils.ConvertStructToMap(request)
		paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
		paymentTransaction.RequestId = strconv.Itoa(responsePay.Data.InvoiceId)
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
		}
		paymentResponse.Transaction = paymentTransaction
		return paymentResponse, validators.ErrorResponse{}
	}

	paymentTransaction.Request = utils.ConvertStructToMap(request)
	paymentTransaction.Response = utils.ConvertStructToMap(responsePay)
	paymentTransaction.Status = consts.PaymentFailedStatus
	paymentTransaction.RequestId = strconv.Itoa(responsePay.Data.InvoiceId)
	errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
	if errUpdate != nil {
		p.logger.Error("Update Error => ", errUpdate)
	}
	paymentResponse.Transaction = paymentTransaction
	return paymentResponse, errRe
}
func UpdateOrderStatus(p PaymentUseCase, ctx context.Context, transaction *domain.Payment) (err validators.ErrorResponse) {
	//err = p.extService.OrderService.SetOrderPaid(ctx, transaction.TransactionId, *transaction)
	//if err.IsError {
	//	p.logger.Error("Unable to update order status => ", err)
	//}
	return
}
