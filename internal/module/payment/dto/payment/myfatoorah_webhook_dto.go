package payment

type MyFatoorahWebhookPayload struct {
	EventType      int    `json:"EventType"`
	Event          string `json:"Event"`
	DateTime       string `json:"DateTime"`
	CountryIsoCode string `json:"CountryIsoCode"`
	Data           struct {
		InvoiceId                    int         `json:"InvoiceId"`
		InvoiceReference             string      `json:"InvoiceReference"`
		CreatedDate                  string      `json:"CreatedDate"`
		CustomerReference            string      `json:"CustomerReference"`
		CustomerName                 string      `json:"CustomerName"`
		CustomerMobile               string      `json:"CustomerMobile"`
		CustomerEmail                string      `json:"CustomerEmail"`
		TransactionStatus            string      `json:"TransactionStatus"`
		PaymentMethod                string      `json:"PaymentMethod"`
		UserDefinedField             interface{} `json:"UserDefinedField"`
		ReferenceId                  string      `json:"ReferenceId"`
		TrackId                      string      `json:"TrackId"`
		PaymentId                    string      `json:"PaymentId"`
		AuthorizationId              string      `json:"AuthorizationId"`
		InvoiceValueInBaseCurrency   string      `json:"InvoiceValueInBaseCurrency"`
		BaseCurrency                 string      `json:"BaseCurrency"`
		InvoiceValueInDisplayCurreny string      `json:"InvoiceValueInDisplayCurreny"`
		DisplayCurrency              string      `json:"DisplayCurrency"`
		InvoiceValueInPayCurrency    string      `json:"InvoiceValueInPayCurrency"`
		PayCurrency                  string      `json:"PayCurrency"`
	} `json:"Data"`
}
