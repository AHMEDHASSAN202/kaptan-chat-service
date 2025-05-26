package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type RejectOffer struct {
	MessageId uint `validate:"required" json:"id" form:"id" query:"message_id" param:"id"`
	dto.MobileHeaders
}

func (input *RejectOffer) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
