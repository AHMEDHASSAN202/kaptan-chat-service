package dto

import (
	"context"
	"github.com/go-playground/validator/v10"
	"kaptan/pkg/utils/dto"
	"kaptan/pkg/validators"
)

type SendMessage struct {
	Channel                 string   `validate:"required" json:"channel" form:"channel" query:"channel"`
	BrandId                 *int64   `json:"brand_id" form:"brand_id" query:"brand_id"`
	Message                 string   `validate:"required" json:"message" form:"message" query:"message"`
	MessageType             string   `validate:"required" json:"message_type" form:"message_type" query:"message_type"`
	TransferId              *int64   `json:"transfer_id" form:"transfer_id" query:"transfer_id"`
	TransferOffersRequested bool     `json:"transfer_offers_requested" form:"transfer_offers_requested" query:"transfer_offers_requested"`
	Price                   *float64 `json:"price" form:"price" query:"price"`
	Note                    *string  `json:"note" form:"note" query:"note"`
	Phone                   *string  `json:"phone" form:"phone" query:"phone"`
	dto.MobileHeaders
}

func (input *SendMessage) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
