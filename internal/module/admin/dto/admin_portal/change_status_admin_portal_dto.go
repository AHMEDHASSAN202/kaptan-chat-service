package admin_portal

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type ChangeAdminPortalStatusDto struct {
	Id           string             `param:"id" validate:"required,mongodb"`
	Status       string             `json:"status" validate:"oneof=active inactive"`
	AdminDetails []dto.AdminDetails `json:"-"`
	dto.PortalHeaders
}

func (input *ChangeAdminPortalStatusDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
