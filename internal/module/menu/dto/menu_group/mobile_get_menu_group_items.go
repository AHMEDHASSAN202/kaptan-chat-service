package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type GetMenuGroupItemDTO struct {
	Query    string `json:"query" form:"query" query:"query"`
	BranchId string `json:"branch_id" form:"branch_id" query:"branch_id" validate:"required,mongodb"`
	dto.MobileHeaders
}

func (input *GetMenuGroupItemDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
