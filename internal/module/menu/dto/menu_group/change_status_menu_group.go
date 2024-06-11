package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ChangeMenuGroupStatusDto struct {
	Id           string             `json:"id" validate:"required,mongodb"`
	Entity       string             `json:"entity" validate:"oneof=menu category item"`
	Status       string             `json:"status" validate:"oneof=active inactive"`
	AdminDetails []dto.AdminDetails `json:"-"`
}

func (input *ChangeMenuGroupStatusDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
