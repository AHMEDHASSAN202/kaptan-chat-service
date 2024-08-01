package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ChangeStatusApprovalDto struct {
	dto.AdminHeaders
	Id     string `param:"id" validate:"required,mongodb"`
	Status string `json:"status" validate:"oneof=rejected approved"`
}

func (input *ChangeStatusApprovalDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
