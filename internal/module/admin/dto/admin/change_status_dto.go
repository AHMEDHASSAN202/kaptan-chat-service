package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ChangeAdminStatusDto struct {
	Id           string             `param:"id" validate:"required,mongodb"`
	Status       string             `json:"status" validate:"oneof=active inactive"`
	AccountId    string             `json:"account_id"`
	AdminDetails []dto.AdminDetails `json:"-"`
}

func (input *ChangeAdminStatusDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
