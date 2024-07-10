package myfatoorah

import (
	"context"
	"github.com/pkg/errors"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/gateways/myfatoorah/consts"
	"samm/internal/module/payment/gateways/myfatoorah/requests"
	"samm/internal/module/payment/gateways/myfatoorah/responses"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

func GetPaymentMethodId(cardType string) int {
	paymentMethods := map[string]int{"mada": 6, "amex": 3, "visa": 2, "master": 2}

	return paymentMethods[cardType]
}
func ExecutePaymentCard(m MyFatoorahService, ctx context.Context, payload *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.ExecutePaymentCardResponse, requestPayload requests.ExecutePaymentCardRequest, err validators.ErrorResponse) {
	//Prepare Request Data
	requestPayload = requests.ExecutePaymentCardRequest{
		PaymentMethodId:    GetPaymentMethodId(payload.Card.Type),
		InvoiceValue:       paymentTransaction.Amount,
		UserDefinedField:   utils.ConvertObjectIdToStringId(paymentTransaction.ID),
		DisplayCurrencyIso: paymentTransaction.Currency,
		CallBackUrl:        m.MainFrontUrl + consts.SuccessUrl,
		ErrorUrl:           m.MainFrontUrl + consts.ErrorUrl,
	}
	requestPayload.ProcessingDetails.AutoCapture = !payload.HoldTransaction
	requestPayload.ProcessingDetails.Bypass3DS = false

	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}
	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&paymentResponse).Post(m.BaseUrl + consts.ExecutePaymentUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> ExecutePaymentCard", errRe)
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> ExecutePaymentCard", errors.New(paymentResponse.Message))
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errors.New(paymentResponse.Message))
	}
	return paymentResponse, requestPayload, err

}

func DirectPaymentCard(m MyFatoorahService, ctx context.Context, url string, payload *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.DirectPaymentResponse, requestPayload requests.DirectPaymentRequest, err validators.ErrorResponse) {
	if payload.PaymentToken != "" {
		//Prepare Request Data
		requestPayload = requests.DirectPaymentRequest{
			PaymentType: "token",
			Bypass3DS:   false,
		}
		requestPayload.Card.SecurityCode = payload.Card.Cvv
		requestPayload.Token = payload.PaymentToken
	} else {
		//Prepare Request Data
		requestPayload = requests.DirectPaymentRequest{
			PaymentType: "card",
			Bypass3DS:   false,
			SaveToken:   payload.SaveCard,
		}
		requestPayload.Card.Number = payload.Card.Number
		requestPayload.Card.HolderName = payload.Card.HolderName
		requestPayload.Card.ExpiryMonth = payload.Card.ExpiryMonth
		requestPayload.Card.ExpiryYear = payload.Card.ExpiryYear
		requestPayload.Card.SecurityCode = payload.Card.Cvv
	}

	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}
	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&paymentResponse).Post(url)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> DirectPaymentCard", errRe)
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> DirectPaymentCard", errors.New(paymentResponse.Message))
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errors.New(paymentResponse.Message))
	}
	return paymentResponse, requestPayload, err

}
func InitSession(ctx context.Context, payload *payment.PayDto, m MyFatoorahService, paymentTransaction *domain.Payment) (initSessionResponse responses.InitSessionResponse, err validators.ErrorResponse) {
	//Prepare Request Data
	requestPayload := requests.InitSessionRequest{
		CustomerIdentifier: utils.ConvertObjectIdToStringId(paymentTransaction.ID),
		SaveToken:          false,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}

	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&initSessionResponse).Post(m.BaseUrl + consts.InitSessionUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> InitSession", errRe)
		return initSessionResponse, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> InitSession", errors.New(initSessionResponse.Message))
		return initSessionResponse, validators.GetErrorResponseFromErr(errors.New(initSessionResponse.Message))
	}
	return initSessionResponse, err
}
func UpdateSession(ctx context.Context, payload *payment.PayDto, sessionId string, m MyFatoorahService) (updateSessionResponse responses.UpdateSessionResponse, err validators.ErrorResponse) {
	requestPayload := requests.UpdateSessionRequest{
		SessionId: sessionId,
		Token:     payload.PaymentToken,
		TokenType: consts.ApplePay,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}

	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&updateSessionResponse).Post(m.BaseUrl + consts.UpdateSessionUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> UpdateSession", errRe)
		return updateSessionResponse, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> UpdateSession", errors.New(updateSessionResponse.Message))
		return updateSessionResponse, validators.GetErrorResponseFromErr(errors.New(updateSessionResponse.Message))
	}
	return updateSessionResponse, err

}
func ExecutePayment(ctx context.Context, payload *payment.PayDto, sessionId string, m MyFatoorahService, paymentTransaction *domain.Payment) (paymentResponse responses.ExecutePaymentResponse, requestPayload requests.ApplePayExecutePaymentCardRequest, err validators.ErrorResponse) {

	//Prepare Request Data
	requestPayload = requests.ApplePayExecutePaymentCardRequest{
		SessionId:          sessionId,
		InvoiceValue:       paymentTransaction.Amount,
		UserDefinedField:   utils.ConvertObjectIdToStringId(paymentTransaction.ID),
		DisplayCurrencyIso: consts.CurrencySAR,
	}
	requestPayload.ProcessingDetails.AutoCapture = !payload.HoldTransaction
	requestPayload.ProcessingDetails.Bypass3DS = true

	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}

	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&paymentResponse).Post(m.BaseUrl + consts.ExecutePaymentUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> ExecutePayment", errRe)
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> ExecutePayment", errors.New(paymentResponse.Message))
		return paymentResponse, requestPayload, validators.GetErrorResponseFromErr(errRe)
	}
	return paymentResponse, requestPayload, validators.ErrorResponse{}
}
