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
