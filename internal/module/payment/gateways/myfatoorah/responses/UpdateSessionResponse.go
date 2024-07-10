package responses

type UpdateSessionResponse struct {
	IsSuccess        bool   `json:"IsSuccess"`
	Message          string `json:"Message"`
	ValidationErrors []struct {
		Name  string `json:"Name"`
		Error string `json:"Error"`
	} `json:"ValidationErrors"`
	Data struct {
		SessionId   string `json:"SessionId"`
		CountryCode string `json:"CountryCode"`
	} `json:"Data"`
}
