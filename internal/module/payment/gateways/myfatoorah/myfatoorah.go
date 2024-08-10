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

func (m MyFatoorahService) InitPaymentCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.InitSessionResponse, redirectUrl string, err validators.ErrorResponse) {

	// call init session with customer With Flag Save Token or not
	paymentResponse, err = InitSession(ctx, dto, m, paymentTransaction)

	if err.IsError {
		return
	}
	redirectUrl = m.MainFrontUrl + "/" + consts.ADD_CARD_VIEW + "?sessionId=" + paymentResponse.Data.SessionId
	return
}
func (m MyFatoorahService) PayCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.ExecutePaymentCardResponse, requestPayload requests.ExecutePaymentCardRequest, invoiceId int, err validators.ErrorResponse) {

	//call init session with customer With Flag Save Token or not

	initSessionResponse, err := InitSession(ctx, dto, m, paymentTransaction)
	//
	if err.IsError {
		return
	}
	updateSessionResponse, err := UpdateSession(ctx, dto, initSessionResponse.Data.SessionId, consts.MFToken, m)

	if err.IsError {
		return
	}
	// Call Execute payment
	executeResponse, requestPayload, err := ExecutePaymentCard(m, ctx, updateSessionResponse.Data.SessionId, paymentTransaction)
	//
	if err.IsError {
		return
	}
	return executeResponse, requestPayload, executeResponse.Data.InvoiceId, err

}
func (m MyFatoorahService) ExecutePaymentCard(ctx context.Context, paymentTransaction *domain.Payment) (paymentResponse responses.ExecutePaymentCardResponse, requestPayload requests.ExecutePaymentCardRequest, invoiceId int, err validators.ErrorResponse) {

	executeResponse, requestPayload, err := ExecutePaymentCard(m, ctx, paymentTransaction.RequestId, paymentTransaction)

	if err.IsError {
		return
	}
	return executeResponse, requestPayload, executeResponse.Data.InvoiceId, err
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

func (m MyFatoorahService) ApplePay(ctx context.Context, dto *payment.PayDto, paymentTransaction *domain.Payment) (paymentResponse responses.ExecutePaymentResponse, requestPayload requests.ApplePayExecutePaymentCardRequest, err validators.ErrorResponse) {

	// Call Init Session
	initResponse, errRe := InitSession(ctx, dto, m, paymentTransaction)
	if errRe.IsError {
		m.logger.Error("Init Session Error ", errRe)
		err = errRe
		return
	}
	// Call Update Session
	updateResponse, errRe := UpdateSession(ctx, dto, initResponse.Data.SessionId, consts.ApplePay, m)
	if errRe.IsError {
		m.logger.Info("Update Session Error ", errRe)
		err = errRe
		return
	}
	// Call Execute Payment
	paymentResponse, paymentRequest, errRe := ExecutePayment(ctx, dto, updateResponse.Data.SessionId, m, paymentTransaction)
	if errRe.IsError {
		m.logger.Info("Execute Payment Error ", errRe)
		err = errRe
		return
	}
	m.httpClient.NewRequest().Get(paymentResponse.Data.PaymentURL)

	return paymentResponse, paymentRequest, validators.ErrorResponse{}
}
func (m MyFatoorahService) GetUserCards(ctx context.Context, userId string) (initSessionResponse responses.InitSessionResponse, err validators.ErrorResponse) {
	requestPayload := requests.InitSessionRequest{
		CustomerIdentifier: userId,
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
func (m MyFatoorahService) DeleteUserCardToken(ctx context.Context, userCardToken string) (err validators.ErrorResponse) {

	headers := map[string]string{
		"Authorization": "Bearer " + m.APIToken,
		"Content-Type":  "application/json",
	}

	res, errRe := m.httpClient.NewRequest().SetHeaders(headers).Post(m.BaseUrl + consts.CancelTokenUrl + "?Token=" + userCardToken)

	if errRe != nil {
		m.logger.Error(ErrorTag+"=> DeleteUserCardToken", errRe)
		return validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return err
}
