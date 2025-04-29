package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type DeleteMessage struct {
	MessageId uint `param:"id"`
	dto.MobileHeaders
}

func (input *DeleteMessage) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
