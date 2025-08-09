package notification

import (
	"context"
	"fmt"
	domain2 "kaptan/internal/module/user/domain"
	"kaptan/pkg/fcm_notification"
	"kaptan/pkg/gate"
	"kaptan/pkg/logger"
)

type UseCase struct {
	logger     logger.ILogger
	gate       *gate.Gate
	driverRepo domain2.DriverRepository
	fcmClient  *fcm_notification.FCMClient
}

func NewNotificationUseCase(driverRepo domain2.DriverRepository, gate *gate.Gate, logger logger.ILogger, fcmClient *fcm_notification.FCMClient) *UseCase {
	return &UseCase{
		logger:     logger,
		gate:       gate,
		driverRepo: driverRepo,
		fcmClient:  fcmClient,
	}
}

func (u *UseCase) SendNotificationToUsers(ctx context.Context, userIDs []uint, title, message string, data map[string]string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Get FCM tokens for users
	tokens, err := u.driverRepo.GetFcmTokenByIds(&ctx, userIDs)
	if err != nil {
		return fmt.Errorf("failed to get FCM tokens: %w", err)
	}

	if len(tokens) == 0 {
		return fmt.Errorf("no FCM tokens found for users")
	}

	// Send batch notification
	batchResponse, err := u.fcmClient.SendToTokens(ctx, tokens, title, message, data)
	if err != nil {
		return fmt.Errorf("failed to send batch notification: %w", err)
	}

	// Handle failed tokens
	var invalidTokens []string
	for _, response := range batchResponse {
		if response.Error != nil {
			if fcm_notification.IsInvalidTokenError(response.Error) {
				invalidTokens = append(invalidTokens, response.Token)
			}
		}
	}

	// Remove invalid tokens from database
	if len(invalidTokens) > 0 {
		_ = u.driverRepo.RemoveInvalidFcmTokens(&ctx, invalidTokens)
	}

	return nil
}

func (u *UseCase) SendNotificationToTopic(ctx context.Context, topic string, title, message string, data map[string]string, condition string) error {
	response, err := u.fcmClient.SendToTopic(ctx, topic, title, message, data, condition)
	if err != nil {
		return fmt.Errorf("failed to send batch notification: %w", err)
	}

	if response.Error != nil {
		return fmt.Errorf("failed to send topic notification: %w", response.Error)
	}

	return nil
}
