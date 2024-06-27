package payment

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type Card struct {
	Type        string `json:"type" validate:"required_if=PaymentType card"`
	Number      string `json:"number" validate:"required_if=PaymentType card"`
	ExpiryMonth string `json:"expiry_month" validate:"required_if=PaymentType card"`
	ExpiryYear  string `json:"expiry_year" validate:"required_if=PaymentType card"`
	Cvv         string `json:"cvv" validate:"required_if=PaymentType card"`
	HolderName  string `json:"holder_name"`
}
type PayDto struct {
	TransactionId   string `json:"transaction_id" validate:"required"`
	TransactionType string `json:"order_type" validate:"required,oneof=order wallet"`
	PaymentType     string `json:"payment_type" validate:"required,oneof=card applepay"`
	PaymentToken    string `json:"payment_token"  validate:"required_if=PaymentType applepay"`
	SaveCard        bool   `json:"save_card"`
	HoldTransaction bool   `json:"hold_transaction"`
	Card            Card   `json:"card" validate:"required_if=PaymentType card"`
}

func (payload *PayDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
