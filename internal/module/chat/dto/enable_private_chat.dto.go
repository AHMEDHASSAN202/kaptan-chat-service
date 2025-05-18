package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type EnablePrivateChat struct {
	Channel string `validate:"required" json:"channel" form:"channel" query:"channel"`
	dto.MobileHeaders
}

func (input *EnablePrivateChat) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
