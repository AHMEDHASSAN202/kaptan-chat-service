package notification

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type StoreNotificationDto struct {
	Title        Name     `json:"title" validate:"required"`
	Text         Name     `json:"Text" validate:"required"`
	Image        string   `json:"image"`
	Type         string   `json:"type" validate:"required,oneof=public private"`
	UserIds      []string `json:"user_ids" validate:"required_if=Type private,max=100"`
	RedirectType string   `json:"redirect_type" validate:"required,oneof=home location"`
	LocationId   string   `json:"location_id" validate:"required_if=RedirectType location"`

	dto.AdminHeaders
}

func (payload *StoreNotificationDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}

type ListNotificationDto struct {
	dto.Pagination
	Query string `query:"query"`
	dto.AdminHeaders
}
