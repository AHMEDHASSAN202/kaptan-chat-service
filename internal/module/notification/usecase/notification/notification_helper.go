package notification

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	domain2 "samm/internal/module/notification/domain"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"strings"
)

func GetNotificationByCode(code string, notificationData map[string]string) (notificationCode domain2.NotificationCode, errResp validators.ErrorResponse) {

	notificationCodes := map[string]domain2.NotificationCode{}

	dir, err := os.Getwd()
	if err != nil {
		return notificationCode, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "notification", "notification_codes", "codes.json")
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error(err)
		return notificationCode, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &notificationCodes); errRe != nil {
		logger.Logger.Error(err)
		return notificationCode, validators.GetErrorResponseFromErr(errRe)
	}
	if notificationToSend, ok := notificationCodes[code]; ok {
		// Replace notification Strings
		for key, value := range notificationData {
			notificationToSend.User.Description.Ar = strings.Replace(notificationToSend.User.Description.Ar, ":"+key, value, 1)
			notificationToSend.User.Description.En = strings.Replace(notificationToSend.User.Description.En, ":"+key, value, 1)
			notificationToSend.Kitchen.Description.Ar = strings.Replace(notificationToSend.Kitchen.Description.Ar, ":"+key, value, 1)
			notificationToSend.Kitchen.Description.En = strings.Replace(notificationToSend.Kitchen.Description.En, ":"+key, value, 1)
		}
		return notificationToSend, validators.ErrorResponse{}
	}
	return notificationCode, validators.GetErrorResponseFromErr(errors.New("Code Not Found"))
}
