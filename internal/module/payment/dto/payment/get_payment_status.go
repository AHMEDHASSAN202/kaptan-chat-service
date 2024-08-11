package payment

import (
	"samm/pkg/utils/dto"
)

type GetPaymentStatus struct {
	TransactionId   string `json:"transaction_id" validate:"required"`
	TransactionType string `json:"transaction_type" validate:"required"`
	UserId          string
	dto.MobileHeaders
}
