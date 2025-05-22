package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type GetChatMessage struct {
	dto.Pagination
	dto.MobileHeaders
	Channel string `validate:"required" json:"channel" form:"channel" query:"channel" param:"channel"`
}

func (input *GetChatMessage) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
