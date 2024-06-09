package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ListMenuGroupDTO struct {
	dto.Pagination
	Query     string `json:"query" form:"query" query:"query"`
	BranchId  string `json:"branch_id" form:"branch_id" query:"branch_id" validate:"omitempty,mongodb"`
	AccountId string `json:"account_id" form:"account_id" validate:"required,mongodb" query:"account_id"`
	Status    string `json:"status" form:"status" validate:"omitempty,oneof=active inactive" query:"status"`
}

func (input *ListMenuGroupDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
