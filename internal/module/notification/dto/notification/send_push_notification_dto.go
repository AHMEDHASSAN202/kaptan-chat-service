package notification

type NotificationDto struct {
	Title        Name     `json:"title"`
	Text         Name     `json:"text"`
	Image        string   `json:"image"`
	Type         string   `json:"type"`
	Ids          []string `json:"ids"`
	AccountIds   []string `json:"account_ids"`
	RedirectType string   `json:"redirect_type"`
	LocationId   string   `json:"location_id"`
	ModelType    string   `json:"model_type"`
	CountryId    string   `json:"country_id"`
}

type NotificationReceiver struct {
	Id              string `json:"id"`
	AccountId       string `json:"account_id"` // in case of location
	Model           string `json:"model"`
	LogNotification bool   `json:"log_notification"` // In case of user only

}
type GeneralNotification struct {
	Country          string                 `json:"country"`
	NotificationData map[string]string      `json:"notification_data"`
	NotificationCode string                 `json:"notification_code"`
	To               []NotificationReceiver `json:"to"`
}
