package requests

type Translated struct {
	Ar string `json:"ar,omitempty"`
	En string `json:"en,omitempty"`
}

type PushNotificationPayload struct {
	AppId            string              `json:"app_id"`
	Headings         Translated          `json:"headings"`
	Contents         Translated          `json:"contents"`
	Name             string              `json:"name"`
	IncludePlayerIds []string            `json:"include_player_ids,omitempty"`
	Filters          []map[string]string `json:"filters,omitempty"`
}
