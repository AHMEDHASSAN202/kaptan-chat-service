package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type UnreadMessages struct {
	Channel *string `json:"channel" form:"channel" query:"channel"`
	dto.MobileHeaders
}

func (input *UnreadMessages) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
