package payment

import (
	"context"
	"fmt"
	"samm/internal/module/payment/consts"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/card"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/external/order/responses"
	"samm/internal/module/payment/response"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strconv"
)

func CreateTransactionBuilder(dto *payment.PayDto, order responses.OrderResponse) *domain.Payment {
	paymentDomain := domain.Payment{}
	paymentDomain.TransactionId = utils.ConvertStringIdToObjectId(dto.TransactionId)
	paymentDomain.TransactionType = dto.TransactionType
	paymentDomain.Amount = order.PriceSummary.TotalPriceAfterDiscount
	paymentDomain.Currency = order.Location.Country.Currency // Todo Depend On Find Order
	paymentDomain.PaymentType = dto.PaymentType

	paymentDomain.Status = consts.PaymentPendingStatus
	// Just for testing
	//paymentDomain.Status = consts.PaymentPaidStatus

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
		hold := len(responsePay.Data.InvoiceTransactions) > 0 && responsePay.Data.InvoiceTransactions[0].TransactionStatus == consts.MyFatoorahHoldInvoiceStatus
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

	// Get Total User Cards
	listCardPayload := card.ListCardDto{
		UserId: dto.UserId,
	}
	listCardPayload.SetDefault()
	userCards, _, errRe := p.cardRepo.ListCard(ctx, &listCardPayload)
	if errRe != nil {
		return response, validators.GetErrorResponseFromErr(errRe)
	}
	if dto.PaymentToken == "" && len(userCards) >= consts.MAX_USER_CARDS {
		return response, validators.GetErrorResponseWithErrors(&ctx, localization.Max_User_Cards, nil)
	}
	// find Order To get Amount
	order, err := p.extService.OrderService.FindOrder(ctx, dto.TransactionId)
	if err.IsError {
		return response, err
	}

	// Create The Transaction
	paymentTransaction, errResp = p.repo.CreateTransaction(ctx, CreateTransactionBuilder(dto, order))
	if errResp != nil {
		return response, validators.GetErrorResponseFromErr(errResp)
	}
	if dto.PaymentToken == "" {
		initSessionResponse, redirectUrl, errS := p.myfatoorahService.InitPaymentCard(ctx, dto, paymentTransaction)
		if errS.IsError {
			p.logger.Error("Init Session Error => ", errS)
			return response, errS
		}
		paymentTransaction.Response = utils.ConvertStructToMap(initSessionResponse)
		paymentTransaction.RequestId = initSessionResponse.Data.SessionId
		errUpdate := p.repo.UpdateTransaction(ctx, paymentTransaction)
		if errUpdate != nil {
			p.logger.Error("Update Error => ", errUpdate)
		}
		response.Transaction = paymentTransaction
		response.RedirectUrl = &redirectUrl
		return
	}
	// Call Myfatoorah

	// Find Card To Get Token
	cardDomain, errRe := p.cardRepo.FindCard(ctx, utils.ConvertStringIdToObjectId(dto.PaymentToken), utils.ConvertStringIdToObjectId(dto.UserId))
	if errRe != nil {
		p.logger.Error("Get Card Error => ", errRe)
		return response, validators.GetErrorResponseFromErr(errRe)
	}
	dto.PaymentToken = cardDomain.MFToken

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

	// find Order To get Amount
	order, err := p.extService.OrderService.FindOrder(ctx, dto.TransactionId)
	if err.IsError {
		return paymentResponse, err
	}
	// Create The Transaction
	paymentTransaction, errResp = p.repo.CreateTransaction(ctx, CreateTransactionBuilder(dto, order))
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
	err = p.extService.OrderService.SetOrderPaid(ctx, transaction.TransactionId, *transaction)
	if err.IsError {
		p.logger.Error("Unable to update order status => ", err)
	}
	return
}
func HandleUpdateUserCards(p PaymentUseCase, ctx context.Context, transaction *domain.Payment) (err validators.ErrorResponse) {

	// Find Order
	order, errRe := p.extService.OrderService.FindOrder(ctx, utils.ConvertObjectIdToStringId(transaction.TransactionId))
	if errRe.IsError {
		fmt.Println("Order Not Found", errRe)
		return errRe
	}

	// call myfatoorah to get user cards
	initSessionResponse, err := p.myfatoorahService.GetUserCards(ctx, utils.ConvertObjectIdToStringId(order.User.ID))
	if err.IsError {

		return err
	}
	fmt.Println("Iam in Update User Cards initSessionResponse ", initSessionResponse)

	// Update User Cards
	userCards := make([]domain.Card, 0)

	for _, cardToken := range initSessionResponse.Data.CustomerTokens {
		userCards = append(userCards, domain.Card{
			Type:    cardToken.CardBrand,
			Number:  cardToken.CardNumber,
			MFToken: cardToken.Token,
			UserId:  order.User.ID,
		})
	}
	er := p.cardRepo.UpdateUserCards(ctx, userCards)

	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	return
}
