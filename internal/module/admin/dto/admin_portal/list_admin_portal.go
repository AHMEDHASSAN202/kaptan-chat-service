package admin_portal

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListAdminPortalDTO struct {
	dto.Pagination
	Query  string `json:"query" form:"query" query:"query"`
	Status string `json:"status" form:"status" validate:"omitempty,oneof=active inactive" query:"status"`
	Role   string `json:"role" form:"role" query:"role"`
	dto.PortalHeaders
}

func (input *ListAdminPortalDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
