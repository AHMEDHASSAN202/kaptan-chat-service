package role

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListRoleDTO struct {
	dto.Pagination
	Query string `json:"query" form:"query" query:"query"`
	dto.AdminHeaders
}

func (input *ListRoleDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
