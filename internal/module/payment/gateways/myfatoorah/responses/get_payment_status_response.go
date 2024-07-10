package responses

type GetPaymentStatusResponse struct {
	IsSuccess        bool        `json:"IsSuccess"`
	Message          string      `json:"Message"`
	ValidationErrors interface{} `json:"ValidationErrors"`
	Data             struct {
		InvoiceId           int           `json:"InvoiceId"`
		InvoiceStatus       string        `json:"InvoiceStatus"`
		InvoiceReference    string        `json:"InvoiceReference"`
		CustomerReference   interface{}   `json:"CustomerReference"`
		CreatedDate         string        `json:"CreatedDate"`
		ExpiryDate          string        `json:"ExpiryDate"`
		ExpiryTime          string        `json:"ExpiryTime"`
		InvoiceValue        float64       `json:"InvoiceValue"`
		Comments            interface{}   `json:"Comments"`
		CustomerName        string        `json:"CustomerName"`
		CustomerMobile      string        `json:"CustomerMobile"`
		CustomerEmail       interface{}   `json:"CustomerEmail"`
		UserDefinedField    string        `json:"UserDefinedField"`
		InvoiceDisplayValue string        `json:"InvoiceDisplayValue"`
		DueDeposit          float64       `json:"DueDeposit"`
		DepositStatus       string        `json:"DepositStatus"`
		InvoiceItems        []interface{} `json:"InvoiceItems"`
		InvoiceTransactions []struct {
			TransactionDate       string      `json:"TransactionDate"`
			PaymentGateway        string      `json:"PaymentGateway"`
			ReferenceId           string      `json:"ReferenceId"`
			TrackId               string      `json:"TrackId"`
			TransactionId         string      `json:"TransactionId"`
			PaymentId             string      `json:"PaymentId"`
			AuthorizationId       string      `json:"AuthorizationId"`
			TransactionStatus     string      `json:"TransactionStatus"`
			TransationValue       string      `json:"TransationValue"`
			CustomerServiceCharge string      `json:"CustomerServiceCharge"`
			TotalServiceCharge    string      `json:"TotalServiceCharge"`
			DueValue              string      `json:"DueValue"`
			PaidCurrency          string      `json:"PaidCurrency"`
			PaidCurrencyValue     string      `json:"PaidCurrencyValue"`
			VatAmount             string      `json:"VatAmount"`
			IpAddress             interface{} `json:"IpAddress"`
			Country               interface{} `json:"Country"`
			Currency              string      `json:"Currency"`
			Error                 string      `json:"Error"`
			CardNumber            interface{} `json:"CardNumber"`
			ErrorCode             string      `json:"ErrorCode"`
		} `json:"InvoiceTransactions"`
		Suppliers []interface{} `json:"Suppliers"`
	} `json:"Data"`
}
