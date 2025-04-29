package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type UpdateMessage struct {
	MessageId uint   `param:"id"`
	Message   string `validate:"required" json:"message" form:"message" query:"message"`
	dto.MobileHeaders
}

func (input *UpdateMessage) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
