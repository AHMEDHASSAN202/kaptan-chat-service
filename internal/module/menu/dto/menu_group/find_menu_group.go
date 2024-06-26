package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type FindMenuGroupDTO struct {
	Id string `param:"id" validate:"required,mongodb"`
	dto.PortalHeaders
}

func (input *FindMenuGroupDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
