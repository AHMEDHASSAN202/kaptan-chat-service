package dto

import (
	"example.com/fxdemo/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"time"
)

type LocationRegisterWebhook struct {
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	UserName  string    `json:"userName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=4"`
	CreatedAt time.Time `validate:"required"`
}

func (l *LocationRegisterWebhook) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c, validate, l)
}
