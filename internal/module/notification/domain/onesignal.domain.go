package domain

import (
	"context"
	"samm/internal/module/notification/gateways/onesignal/requests"
	"samm/internal/module/notification/gateways/onesignal/responses"
	"samm/pkg/validators"
)

type OneSignalService interface {
	SendNotification(ctx context.Context, dto *requests.PushNotificationPayload, modelType string) (notificationResponse responses.PushNotificationResponse, err validators.ErrorResponse)
}
