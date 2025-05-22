package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
	"strings"
)

type GetChat struct {
	Channel  string `validate:"required" json:"channel" form:"channel" query:"channel" param:"channel"`
	MarkRead string `form:"mark_read" query:"mark_read"`
	dto.MobileHeaders
}

func (input *GetChat) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}

func (input *GetChat) GetMarkAsRead() bool {
	return strings.ToLower(input.MarkRead) == "true"
}
