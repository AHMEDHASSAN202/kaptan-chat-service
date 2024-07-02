package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListUserDto struct {
	dto.Pagination
	Query  string `query:"query"`
	Status string `query:"status" validate:"omitempty,oneof=active inactive"`
	Dob    string `query:"dob" validate:"omitempty,datetime=2006-01-02"`
}

func (input *ListUserDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
