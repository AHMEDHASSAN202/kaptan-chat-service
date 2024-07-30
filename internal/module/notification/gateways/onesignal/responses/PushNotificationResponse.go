package responses

type PushNotificationResponse struct {
	ID     string      `json:"id"`
	Errors interface{} `json:"errors"`
}
