package requests

type UpdatePaymentRequest struct {
	Operation string  `json:"Operation"`
	Amount    float64 `json:"Amount"`
	Key       string  `json:"Key"`
	KeyType   string  `json:"KeyType"`
}
