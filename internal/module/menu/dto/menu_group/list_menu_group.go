package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListMenuGroupDTO struct {
	dto.Pagination
	Query      string `json:"query" form:"query" query:"query"`
	LocationId string `json:"location_id" form:"location_id" query:"location_id" validate:"omitempty,mongodb"`
	Status     string `json:"status" form:"status" validate:"omitempty,oneof=active inactive" query:"status"`
	dto.PortalHeaders
}

func (input *ListMenuGroupDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
