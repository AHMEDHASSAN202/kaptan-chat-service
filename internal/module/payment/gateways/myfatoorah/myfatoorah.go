package myfatoorah

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"os"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
	"samm/internal/module/payment/gateways/myfatoorah/consts"
	"samm/internal/module/payment/gateways/myfatoorah/requests"
	"samm/internal/module/payment/gateways/myfatoorah/responses"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

const ErrorTag = "MyFatoorahService"

type MyFatoorahService struct {
	BaseUrl      string
	APIToken     string
	MainAppUrl   string
	MainFrontUrl string
	logger       logger.ILogger
	httpClient   *resty.Client
}

func NewMyFatoorahService(httpClient *resty.Client, logger logger.ILogger) domain.MyFatoorahService {

	return &MyFatoorahService{
		httpClient:   httpClient,
		logger:       logger,
		BaseUrl:      os.Getenv("MY_FATOORAH_URL"),
		APIToken:     os.Getenv("MY_FATOORAH_SECRET"),
		MainAppUrl:   os.Getenv("MAIN_APP_URL"),
		MainFrontUrl: os.Getenv("MAIN_FRONT_URL"),
	}
}

func (m MyFatoorahService) PayCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.DirectPaymentResponse, requestPayload requests.DirectPaymentRequest, invoiceId int, err validators.ErrorResponse) {

	// Call Execute payment
	executeResponse, _, err := ExecutePaymentCard(m, ctx, dto, paymentTransaction)

	if err.IsError {
		return
	}
	if executeResponse.Data.IsDirectPayment {
		// Call Direct Payment
		directResponse, directRequest, err := DirectPaymentCard(m, ctx, executeResponse.Data.PaymentURL, dto, paymentTransaction)
		return directResponse, directRequest, executeResponse.Data.InvoiceId, err
	}
	return paymentResponse, requestPayload, executeResponse.Data.InvoiceId, validators.GetErrorResponseFromErr(errors.New("No Direct Payment"))

}

func (m MyFatoorahService) FindPayment(ctx context.Context, invoiceId string) (paymentResponse responses.GetPaymentStatusResponse, err validators.ErrorResponse) {

	//Prepare Request Data
	requestPayload := requests.GetPaymentStatusRequest{
		Key:     invoiceId,
		KeyType: consts.KeyTypeInvoiceId,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}
	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(requestPayload).SetResult(&paymentResponse).Post(m.BaseUrl + consts.GetPaymentStatusUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> GetPaymentStatus", errRe)
		return paymentResponse, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> GetPaymentStatus", errors.New(paymentResponse.Message))
		return paymentResponse, validators.GetErrorResponseFromErr(errors.New(paymentResponse.Message))
	}
	return paymentResponse, err
}

func (m MyFatoorahService) UpdatePaymentStatus(ctx context.Context, invoiceId string, capture bool) (err validators.ErrorResponse) {
	responsePay, err := m.FindPayment(ctx, invoiceId)
	if err.IsError {
		return
	}
	operation := "Release"
	if capture {
		operation = "Capture"
	}
	paymentRequest := requests.UpdatePaymentRequest{
		Operation: operation, //Capture
		Amount:    responsePay.Data.InvoiceValue,
		Key:       invoiceId,
		KeyType:   "InvoiceId",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}
	var paymentResponse responses.UpdatePaymentResponse
	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).SetBody(paymentRequest).SetResult(&paymentResponse).Post(m.BaseUrl + consts.UpdatePaymentUrl)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> UpdatePaymentStatus", errRe)
		return validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		m.logger.Error(ErrorTag+"=> UpdatePaymentStatus", errors.New(paymentResponse.Message))
		return validators.GetErrorResponseFromErr(errors.New(paymentResponse.Message))
	}
	return err
}
