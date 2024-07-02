package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UserSignUpDto struct {
	//PhoneNumber string `json:"phone_number" validate:"required,PhoneNumber_rule_validation"`
	//CountryCode string `json:"country_code" validate:"required,len=4,numeric"`
	Name string `json:"name" validate:"required"`
	dto.MobileHeaders
	//Email       string `json:"email" validate:"required,email"`
}

func (payload *UserSignUpDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
