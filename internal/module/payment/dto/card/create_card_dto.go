package card

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type CreateCardDto struct {
	Type        string `json:"type" validate:"required,oneof=visa master mada"`
	Number      string `json:"number" validate:"required"`
	ExpiryMonth string `json:"expiry_month" validate:"required"`
	ExpiryYear  string `json:"expiry_year" validate:"required"`
	Cvv         string `json:"cvv" validate:"required"`
	HolderName  string `json:"holder_name" `

	UserId string
}

func (payload *CreateCardDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
