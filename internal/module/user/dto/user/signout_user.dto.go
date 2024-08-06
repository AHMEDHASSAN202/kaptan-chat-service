package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UserSignOutDto struct {
	Authorization string `header:"Authorization" validate:"required"`
	dto.MobileHeaders
}

func (payload *UserSignOutDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
