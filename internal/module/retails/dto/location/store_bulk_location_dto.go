package location

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type LocationDto struct {
	Name            Name           `json:"name" validate:"required"`
	City            City           `json:"city" validate:"required"`
	Street          Name           `json:"street" validate:"required"`
	Tags            []string       `json:"tags" validate:"required"`
	CoverImage      string         `json:"cover_image"`
	Logo            string         `json:"logo" `
	Phone           string         `json:"phone" validate:"required"`
	Lat             float64        `json:"lat" validate:"required,latitude"`
	Lng             float64        `json:"lng" validate:"required,longitude"`
	WorkingHour     []WorkingHour  `json:"working_hour" `
	Percent         float64        ` json:"percent" validate:"required"`
	PercentsDate    []PercentsDate `json:"percents_date"`
	PreparationTime int            `json:"preparation_time" validate:"required"`
	AutoAccept      bool           `json:"auto_accept" `
	BankAccount     BankAccount    `json:"bank_account" validate:"required"`
}
type StoreBulkLocationDto struct {
	BrandDetails Brand         `json:"brand_details" validate:"required" `
	Country      Country       `json:"country" validate:"required"`
	AccountId    string        `json:"account_id" validate:"required,mongodb"`
	Locations    []LocationDto `json:"locations" validate:"required" `
	dto.AdminHeaders
}

func (payload *StoreBulkLocationDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
