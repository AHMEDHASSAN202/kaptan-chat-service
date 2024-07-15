package order

type CancelOrderDto struct {
	OrderId        string
	UserId         string `header:"causer-id" validate:"required"`
	CancelReasonId string `json:"cancel_reason_id" validate:"required"`
}
