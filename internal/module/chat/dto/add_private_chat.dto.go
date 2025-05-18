package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type AddPrivateChat struct {
	MessageId uint `validate:"required" json:"message_id" form:"message_id" query:"message_id"`
	dto.MobileHeaders
}

func (input *AddPrivateChat) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
