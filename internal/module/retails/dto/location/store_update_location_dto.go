package location

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type Country struct {
	Id   string `json:"_id" validate:"required"`
	Name struct {
		Ar string `json:"ar" validate:"required"`
		En string `json:"en" validate:"required"`
	} `json:"name" validate:"required"`
	Timezone    string `json:"timezone" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
}

type PercentsDate struct {
	From    time.Time `json:"from" validate:"datetime"`
	To      time.Time `json:"to" validate:"datetime"`
	Percent float64   ` json:"percent" validate:"number"`
}

type City struct {
	Id   string `json:"_id" validate:"required,mongodb"`
	Name Name   `json:"name" validate:"required"`
}
type Brand struct {
	Id       string `json:"_id" validate:"required,mongodb"`
	Name     Name   `json:"name" validate:"required" `
	Logo     string `json:"logo" `
	IsActive bool   `json:"is_active"`
}
type WorkingHour struct {
	Day       string `json:"day"  validate:"required"`
	From      string `json:"from" validate:"required_if=IsFullDay false,Timeformat"`
	To        string `json:"to" validate:"required_if=IsFullDay false,Timeformat"`
	IsFullDay bool   `json:"is_full_day" `
}
type BankAccount struct {
	AccountNumber string `json:"account_number" validate:"required"`
	BankName      string `json:"bank_name" validate:"required"`
	CompanyName   string `json:"company_name" validate:"required"`
}

type StoreLocationDto struct {
	Name                       Name           `json:"name" validate:"required"`
	City                       City           `json:"city" validate:"required"`
	Street                     Name           `json:"street" validate:"required"`
	Tags                       []string       `json:"tags" validate:"required"`
	CoverImage                 string         `json:"cover_image"`
	Logo                       string         `json:"logo"`
	Phone                      string         `json:"phone" validate:"required"`
	Lat                        float64        `json:"lat" validate:"required,latitude"`
	Lng                        float64        `json:"lng" validate:"required,longitude"`
	BrandDetails               Brand          `json:"brand_details" validate:"required" `
	WorkingHour                []WorkingHour  `json:"working_hour" `
	WorkingHourEid             []WorkingHour  `json:"working_hour_eid"`
	WorkingHourRamadan         []WorkingHour  `json:"working_hour_ramadan"`
	Percent                    float64        `json:"percent" `
	PercentsDate               []PercentsDate `json:"percents_date"`
	PreparationTime            int            `json:"preparation_time" validate:"required"`
	AutoAccept                 bool           `json:"auto_accept" `
	BankAccount                BankAccount    `json:"bank_account" validate:"required"`
	AccountId                  string         `json:"account_id" validate:"required,mongodb"`
	Country                    Country        `json:"country" validate:"required"`
	AllowedCollectionMethodIds []string       `json:"allowed_collection_method_ids" validate:"required,min=1"`
	dto.AdminHeaders
}

func (payload *StoreLocationDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, payload)
}
