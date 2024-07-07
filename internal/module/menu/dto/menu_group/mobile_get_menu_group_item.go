package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type GetMenuGroupItemDTO struct {
	ID         string `param:"id" validate:"required,mongodb"`
	LocationId string `param:"location_id" validate:"required,mongodb"`
	dto.MobileHeaders
}

func (input *GetMenuGroupItemDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
