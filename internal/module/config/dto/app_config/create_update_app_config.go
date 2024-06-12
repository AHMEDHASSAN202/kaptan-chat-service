package app_config

import (
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CreateUpdateAppConfigDto struct {
	Id           string             `json:"_"`
	ForceUpdate  bool               `json:"force_update"`
	Type         string             `json:"type" validate:"required,oneof=user merchant"`
	StartupImage string             `json:"stratup_image"`
	AdminDetails []dto.AdminDetails `json:"-"`
}

func (input *CreateUpdateAppConfigDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
