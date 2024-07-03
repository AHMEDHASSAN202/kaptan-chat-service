package payment

type AuthorizePayload struct {
	Capture       bool   `json:"capture" form:"capture"`
	TransactionId string `json:"transaction_id" form:"transaction_id"`
}
