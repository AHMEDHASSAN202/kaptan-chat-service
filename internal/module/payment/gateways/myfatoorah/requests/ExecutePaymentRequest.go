package requests

type ExecutePaymentCardRequest struct {
	PaymentMethodId    int     `json:"PaymentMethodId"`
	DisplayCurrencyIso string  `json:"DisplayCurrencyIso"`
	InvoiceValue       float64 `json:"invoiceValue"`
	ProcessingDetails  struct {
		AutoCapture bool `json:"AutoCapture"`
		Bypass3DS   bool `json:"Bypass3DS"`
	} `json:"ProcessingDetails"`
	UserDefinedField string `json:"UserDefinedField"`
	CallBackUrl      string `json:"CallBackUrl"`
	ErrorUrl         string `json:"ErrorUrl"`
}
