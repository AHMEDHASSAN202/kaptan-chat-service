package notification

type NotificationDto struct {
	Title        Name     `json:"title"`
	Text         Name     `json:"text"`
	Image        string   `json:"image"`
	Type         string   `json:"type"`
	Ids          []string `json:"ids"`
	RedirectType string   `json:"redirect_type"`
	LocationId   string   `json:"location_id"`
	ModelType    string   `json:"model_type"`
}
