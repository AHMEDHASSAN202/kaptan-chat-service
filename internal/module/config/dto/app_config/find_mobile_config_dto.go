package app_config

import (
	"context"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
)

type FindMobileConfigDto struct {
	Type     string `json:"type" form:"type" query:"type" validate:"required,oneof=user merchant"`
	Platform string `json:"platform" form:"platform" query:"platform" validate:"required,oneof=ios android huawei"`
	Version  int64  `json:"version" form:"version" query:"version" validate:"required"`
}

func (input *FindMobileConfigDto) Validate(ctx context.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(ctx, validate, input)
}
