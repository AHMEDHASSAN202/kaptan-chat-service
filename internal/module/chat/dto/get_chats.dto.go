package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type GetChats struct {
	dto.MobileHeaders
}

func (input *GetChats) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
