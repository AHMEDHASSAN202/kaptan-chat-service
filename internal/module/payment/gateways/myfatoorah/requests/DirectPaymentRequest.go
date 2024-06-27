package requests

type DirectPaymentRequest struct {
	PaymentType string `json:"PaymentType"`
	Card        struct {
		Number       string `json:"Number"`
		ExpiryMonth  string `json:"ExpiryMonth"`
		ExpiryYear   string `json:"ExpiryYear"`
		SecurityCode string `json:"SecurityCode"`
		HolderName   string `json:"HolderName"`
	} `json:"Card"`

	Token     string `json:"Token"`
	Bypass3DS bool   `json:"Bypass3DS"`
}
