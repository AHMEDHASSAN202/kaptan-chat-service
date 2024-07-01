package responses

import "time"

type UpdatePaymentResponse struct {
	IsSuccess        bool   `json:"IsSuccess"`
	Message          string `json:"Message"`
	ValidationErrors []struct {
		Name  string `json:"Name"`
		Error string `json:"Error"`
	} `json:"ValidationErrors"`
	Data struct {
		InvoiceId           int        `json:"InvoiceId"`
		InvoiceStatus       string     `json:"InvoiceStatus"`
		InvoiceReference    string     `json:"InvoiceReference"`
		CustomerReference   string     `json:"CustomerReference"`
		CreatedDate         *time.Time `json:"CreatedDate"`
		ExpiryDate          string     `json:"ExpiryDate"`
		ExpiryTime          string     `json:"ExpiryTime"`
		InvoiceValue        int        `json:"InvoiceValue"`
		Comments            string     `json:"Comments"`
		CustomerName        string     `json:"CustomerName"`
		CustomerMobile      string     `json:"CustomerMobile"`
		CustomerEmail       string     `json:"CustomerEmail"`
		UserDefinedField    string     `json:"UserDefinedField"`
		InvoiceDisplayValue string     `json:"InvoiceDisplayValue"`
		DueDeposit          int        `json:"DueDeposit"`
		DepositStatus       string     `json:"DepositStatus"`
		InvoiceItems        []struct {
			ItemName  string `json:"ItemName"`
			Quantity  int    `json:"Quantity"`
			UnitPrice int    `json:"UnitPrice"`
			Weight    int    `json:"Weight"`
			Width     int    `json:"Width"`
			Height    int    `json:"Height"`
			Depth     int    `json:"Depth"`
		} `json:"InvoiceItems"`
		InvoiceTransactions []struct {
			TransactionDate       *time.Time `json:"TransactionDate"`
			PaymentGateway        string     `json:"PaymentGateway"`
			ReferenceId           string     `json:"ReferenceId"`
			TrackId               string     `json:"TrackId"`
			TransactionId         string     `json:"TransactionId"`
			PaymentId             string     `json:"PaymentId"`
			AuthorizationId       string     `json:"AuthorizationId"`
			TransactionStatus     string     `json:"TransactionStatus"`
			TransationValue       string     `json:"TransationValue"`
			CustomerServiceCharge string     `json:"CustomerServiceCharge"`
			TotalServiceCharge    string     `json:"TotalServiceCharge"`
			DueValue              string     `json:"DueValue"`
			PaidCurrency          string     `json:"PaidCurrency"`
			PaidCurrencyValue     string     `json:"PaidCurrencyValue"`
			VatAmount             string     `json:"VatAmount"`
			IpAddress             string     `json:"IpAddress"`
			Country               string     `json:"Country"`
			Currency              string     `json:"Currency"`
			Error                 string     `json:"Error"`
			CardNumber            string     `json:"CardNumber"`
			ErrorCode             string     `json:"ErrorCode"`
		} `json:"InvoiceTransactions"`
		Suppliers []struct {
			SupplierCode  int    `json:"SupplierCode"`
			SupplierName  string `json:"SupplierName"`
			InvoiceShare  int    `json:"InvoiceShare"`
			ProposedShare int    `json:"ProposedShare"`
			DepositShare  int    `json:"DepositShare"`
		} `json:"Suppliers"`
	} `json:"Data"`
}
