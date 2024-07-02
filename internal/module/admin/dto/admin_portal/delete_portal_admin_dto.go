package admin_portal

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type DeletePortalAdminDTO struct {
	ID string `param:"id" validate:"required,mongodb"`
	dto.PortalHeaders
}

func (input *DeletePortalAdminDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStructAndReturnOneError(c.Request().Context(), validate, input)
}
