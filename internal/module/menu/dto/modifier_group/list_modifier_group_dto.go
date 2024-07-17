package modifier_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListModifierGroupsDto struct {
	dto.Pagination
	Query     string `json:"query" form:"query" query:"query"`
	AccountId string `json:"account_id" query:"account_id" header:"account_id" validate:"required"`
}

func (input *ListModifierGroupsDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
