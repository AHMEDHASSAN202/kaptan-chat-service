package responses

type DirectPaymentResponse struct {
	IsSuccess        bool   `json:"IsSuccess"`
	Message          string `json:"Message"`
	ValidationErrors []struct {
		Name  string `json:"Name"`
		Error string `json:"Error"`
	} `json:"ValidationErrors"`
	Data struct {
		Status       string `json:"status"`
		ErrorMessage string `json:"ErrorMessage"`
		PaymentId    string `json:"PaymentId"`
		Token        string `json:"Token"`
		PaymentURL   string `json:"PaymentURL"`
	} `json:"Data"`
}
