package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type GetChats struct {
	MessageId  *uint `json:"message_id" form:"message_id" query:"message_id"`
	TransferId *uint `json:"transfer_id" form:"transfer_id" query:"transfer_id"`
	dto.MobileHeaders
}

func (input *GetChats) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
