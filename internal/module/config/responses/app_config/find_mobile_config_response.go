package app_config

type Phone struct {
	PhoneNumber    string `json:"phone_num" bson:"phone_num"`
	WhatsappNumber string `json:"whatsapp_number" bson:"whatsapp_number"`
}
type FindMobileConfigResponse struct {
	ForceUpdate         bool   `json:"force_update"`
	Type                string `json:"type"`
	MinVersion          int64  `json:"min_version"`
	AppLink             string `json:"app_link"`
	LocalizationVersion int64  `json:"localization_version"`
	Phone               Phone  `json:"phone" bson:"phone"`
	StartupImage        string `json:"stratup_image"`
}
