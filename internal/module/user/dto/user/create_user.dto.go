package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
	"time"
)

type CreateUserDto struct {
	Name          string    `json:"name" validate:"required"`
	CountryCode   string    `json:"country_code" validate:"required,len=2"`
	PhoneNumber   string    `json:"phone_number" validate:"required,e164"`
	Email         string    `json:"email" validate:"required,email"`
	Gender        string    `json:"gender" validate:"required,oneof=male female other"`
	Dob           time.Time `json:"dob" validate:"required"`
	Otp           string    `json:"otp" validate:"required,len=6,numeric"`
	ExpiryOtpDate time.Time `json:"expiry_otp_date" validate:"required"`
	OtpCounter    int       `json:"otp_counter" validate:"required,min=0,max=24"`
	ImageURL      string    `json:"image_url" validate:"omitempty,url"`
	Country       string    `json:"country" validate:"required"`
	IsActive      bool      `json:"is_active" validate:"required"`
	DeletedAt     time.Time `json:"deleted_at" validate:"omitempty"`
	Tokens        []string  `json:"tokens" validate:"dive,required"`
	CreatedAt     time.Time `json:"createdAt" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt" validate:"required"`
}

func (payload *CreateUserDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
