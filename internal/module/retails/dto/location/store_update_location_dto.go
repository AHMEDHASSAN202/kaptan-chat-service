package location

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type City struct {
	Id   string `json:"_id" validate:"required,mongodb"`
	Name Name   `json:"name" validate:"required"`
}
type Brand struct {
	Id       string `json:"_id" validate:"required,mongodb"`
	Name     Name   `json:"name" validate:"required" `
	Logo     string `json:"logo" validate:"required"`
	IsActive bool   `json:"is_active"`
}
type WorkingHour struct {
	Day  string `json:"day"  validate:"required"`
	From string `json:"from" validate:"required,timeformat"`
	To   string `json:"to" validate:"required,timeformat"`
}
type BankAccount struct {
	AccountNumber string `json:"account_number" validate:"required"`
	BankName      string `json:"bank_name" validate:"required"`
	CompanyName   string `json:"company_name" validate:"required"`
}

type StoreLocationDto struct {
	Name            Name          `json:"name" validate:"required"`
	City            City          `json:"city" validate:"required"`
	Street          Name          `json:"street" validate:"required"`
	Tags            string        `json:"tags" validate:"required"`
	CoverImage      string        `json:"cover_image"`
	Logo            string        `json:"logo" `
	Phone           string        `json:"phone" validate:"required"`
	Lat             float64       `json:"lat" validate:"required,latitude"`
	Lng             float64       `json:"lng" validate:"required,longitude"`
	BrandDetails    Brand         `json:"brand_details" validate:"required" `
	WorkingHour     []WorkingHour `json:"working_hour" validate:"required"`
	PreparationTime int           `json:"preparation_time" validate:"required"`
	AutoAccept      bool          `json:"auto_accept" `
	BankAccount     BankAccount   `json:"bank_account" validate:"required"`
	AccountId       string        `json:"account_id" validate:"required,mongodb"`
}

func (payload *StoreLocationDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
