package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"time"
)

type UpdateUserProfileDto struct {
	ID          string    `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	CountryCode string    `json:"country_code" validate:"required,len=2"`
	PhoneNumber string    `json:"phone_number" validate:"required,e164"`
	Email       string    `json:"email" validate:"required,email"`
	Gender      string    `json:"gender" validate:"required,oneof=male female other"`
	Dob         time.Time `json:"dob" validate:"required"`
	ImageURL    string    `json:"image_url" validate:"omitempty,url"`
}

func (payload *UpdateUserProfileDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
