package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DeleteEntityFromMenuGroupDto struct {
	Id           string             `param:"id" validate:"required,mongodb"`
	EntityId     string             `param:"entity_id" validate:"required,mongodb"`
	Entity       string             `param:"entity" validate:"oneof=category item"`
	AdminDetails []dto.AdminDetails `json:"-"`
}

func (input *DeleteEntityFromMenuGroupDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
