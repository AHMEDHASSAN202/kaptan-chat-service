package requests

type ApplePayExecutePaymentCardRequest struct {
	SessionId          string  `json:"SessionId"`
	InvoiceValue       float64 `json:"invoiceValue"`
	UserDefinedField   string  `json:"UserDefinedField"`
	DisplayCurrencyIso string  `json:"DisplayCurrencyIso"`
	ProcessingDetails  struct {
		AutoCapture bool `json:"AutoCapture"`
		Bypass3DS   bool `json:"Bypass3DS"`
	} `json:"ProcessingDetails"`
}
