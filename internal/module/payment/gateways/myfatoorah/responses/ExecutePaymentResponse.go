package responses

type ExecutePaymentResponse struct {
	IsSuccess        bool   `json:"IsSuccess"`
	Message          string `json:"Message"`
	ValidationErrors []struct {
		Name  string `json:"Name"`
		Error string `json:"Error"`
	} `json:"ValidationErrors"`
	Data struct {
		InvoiceId         int    `json:"InvoiceId"`
		IsDirectPayment   bool   `json:"IsDirectPayment"`
		PaymentURL        string `json:"PaymentURL"`
		CustomerReference string `json:"CustomerReference"`
		UserDefinedField  string `json:"UserDefinedField"`
		RecurringId       string `json:"RecurringId"`
	} `json:"Data"`
}
