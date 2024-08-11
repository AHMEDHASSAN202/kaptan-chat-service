package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DeleteUserDto struct {
	DeleteReasonId string `json:"reason_id" validate:"required"`
	dto.MobileHeaders
}

func (d *DeleteUserDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, d)
}
