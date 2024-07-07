package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type VerifyUserOtpDto struct {
	PhoneNumber string `json:"phone_number" validate:"required,PhoneNumber_rule_validation"`
	CountryCode string `json:"country_code" validate:"required,oneof=+966 +20 +971"`
	Otp         string `json:"otp" validate:"required,len=4,alphanum"`
}

func (payload *VerifyUserOtpDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
