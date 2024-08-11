package requests

type ApplePayData struct {
	PaymentData struct {
		Data   string `json:"data"`
		Header struct {
			EphemeralPublicKey string `json:"ephemeralPublicKey"`
			PublicKeyHash      string `json:"publicKeyHash"`
			TransactionId      string `json:"transactionId"`
		} `json:"header"`
		Signature string `json:"signature"`
		Version   string `json:"version"`
	} `json:"PaymentData"`
	PaymentMethod struct {
		DisplayName string `json:"displayName"`
		Network     string `json:"network"`
		Type        string `json:"type"`
	} `json:"PaymentMethod"`
	TransactionIdentifier string `json:"TransactionIdentifier"`
}
