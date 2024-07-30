package onesignal

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"os"
	"samm/internal/module/notification/domain"
	"samm/internal/module/notification/gateways/onesignal/consts"
	"samm/internal/module/notification/gateways/onesignal/requests"
	"samm/internal/module/notification/gateways/onesignal/responses"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

const ErrorTag = "OnesignalService"

type OnesignalService struct {
	BaseUrl         string
	APITokenUser    string
	APITokenKitchen string
	UserAppId       string
	KitchenAppId    string
	logger          logger.ILogger
	httpClient      *resty.Client
}

func NewOnesignalService(httpClient *resty.Client, logger logger.ILogger) domain.OneSignalService {

	return &OnesignalService{
		httpClient:      httpClient,
		logger:          logger,
		BaseUrl:         os.Getenv("ONE_SIGNAL_URL"),
		APITokenUser:    os.Getenv("ONE_SIGNAL_SECRET_USER"),
		APITokenKitchen: os.Getenv("ONE_SIGNAL_SECRET_KITCHEN"),
		UserAppId:       os.Getenv("ONE_SIGNAL_APP_ID_USER"),
		KitchenAppId:    os.Getenv("ONE_SIGNAL_APP_ID_KITCHEN"),
	}
}
func (o OnesignalService) GetValidCred(modelType string) (appId, apiToken string) {
	if modelType == consts.KitchenModelType {
		return o.KitchenAppId, o.APITokenKitchen
	}
	return o.UserAppId, o.APITokenUser

}
func (o OnesignalService) SendNotification(ctx context.Context, dto *requests.PushNotificationPayload, modelType string) (notificationResponse responses.PushNotificationResponse, err validators.ErrorResponse) {
	appId, apiToken := o.GetValidCred(modelType)
	dto.AppId = appId
	headers := map[string]string{
		"Authorization": "Basic " + apiToken,
		"Content-Type":  "application/json",
	}
	responseBody := map[string]interface{}{}
	res, errRe := o.httpClient.NewRequest().SetDebug(true).SetHeaders(headers).SetBody(dto).SetResult(&responseBody).Post(o.BaseUrl + consts.SendNotificationsUrl)

	if errRe != nil {
		o.logger.Error(ErrorTag+"=> SendNotification", errRe)
		return notificationResponse, validators.GetErrorResponseFromErr(errRe)
	}
	if !res.IsSuccess() {
		o.logger.Error(ErrorTag+"=> SendNotification", notificationResponse)
		return notificationResponse, validators.GetErrorResponseFromErr(errors.New("Notification Fails"))
	}
	return notificationResponse, err

}
