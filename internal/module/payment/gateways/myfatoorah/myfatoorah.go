package myfatoorah

import (
	"context"
	"github.com/go-resty/resty/v2"
	"os"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/payment"
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

func (m MyFatoorahService) PayCard(ctx context.Context, dto *payment.PayDto, paymentTransaction *domain.Payment) (err validators.ErrorResponse) {

	// Call Execute payment
	//executeResponse, _, err := ExecutePaymentCard(m, ctx, dto, paymentTransaction)

	if err.IsError {
		return
	}
	// Call Direct Payment
	//executeResponse.Data.PaymentURL

	return err
	panic("implement me")
}
