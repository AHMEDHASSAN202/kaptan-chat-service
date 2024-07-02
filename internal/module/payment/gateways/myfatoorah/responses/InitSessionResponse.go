package responses

type InitSessionResponse struct {
	IsSuccess        bool   `json:"IsSuccess"`
	Message          string `json:"Message"`
	ValidationErrors []struct {
		Name  string `json:"Name"`
		Error string `json:"Error"`
	} `json:"ValidationErrors"`
	Data struct {
		SessionId      string `json:"SessionId"`
		CountryCode    string `json:"CountryCode"`
		CustomerTokens []struct {
			Token      string `json:"Token"`
			CardNumber string `json:"CardNumber"`
			CardBrand  string `json:"CardBrand"`
		} `json:"CustomerTokens"`
	} `json:"Data"`
}
