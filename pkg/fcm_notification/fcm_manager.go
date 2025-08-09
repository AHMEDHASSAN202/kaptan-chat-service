package fcm_notification

import (
	"fmt"
	"kaptan/pkg/config"
	"kaptan/pkg/logger"
	"os"
)

func NewFcmManager(log logger.ILogger, config *config.Config) *FCMClient {
	wd, _ := os.Getwd()
	firebaseConfigFile := fmt.Sprintf("%s/%s", wd, config.FirebaseConfig.FcmFilePath)
	fcmClient, err := NewFCMClientFromFile(firebaseConfigFile)
	if err != nil {
		log.Fatalf("Failed to initialize FCM client: %v", err)
	}
	return fcmClient
}
