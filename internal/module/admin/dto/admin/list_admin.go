package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListAdminDTO struct {
	dto.Pagination
	Query  string `json:"query" form:"query" query:"query"`
	Status string `json:"status" form:"status" validate:"omitempty,oneof=active inactive" query:"status"`
	Type   string `json:"type" form:"type" query:"type"`
	Role   string `json:"role" form:"role" query:"role"`
	dto.AdminHeaders
}

func (input *ListAdminDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
