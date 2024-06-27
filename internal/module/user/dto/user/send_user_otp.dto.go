package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/validators"
)

type SendUserOtpDto struct {
	PhoneNumber string `json:"phone_number" validate:"required,PhoneNumber_rule_validation"`
	CountryCode string `json:"country_code" validate:"required,len=4,numeric"`
}

func (payload *SendUserOtpDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
