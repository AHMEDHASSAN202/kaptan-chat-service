package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListApprovalDto struct {
	dto.AdminHeaders
	dto.Pagination
	Type string `json:"type" query:"type" validate:"required"`
}

func (input *ListApprovalDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
