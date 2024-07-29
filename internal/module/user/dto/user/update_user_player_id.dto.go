package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UpdateUserPlayerId struct {
	PlayerId string `json:"player_id" validate:"required"`
	dto.MobileHeaders
}

func (payload *UpdateUserPlayerId) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, payload)
}
