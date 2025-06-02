package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type EndTransfer struct {
	TransferId uint `validate:"required" json:"transfer_id" form:"transfer_id" query:"transfer-id" param:"id"`
	dto.MobileHeaders
}

func (input *EndTransfer) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
