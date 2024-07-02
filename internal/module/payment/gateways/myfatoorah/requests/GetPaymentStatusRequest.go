package requests

type GetPaymentStatusRequest struct {
	Key     string `json:"Key"`
	KeyType string `json:"KeyType"`
}
