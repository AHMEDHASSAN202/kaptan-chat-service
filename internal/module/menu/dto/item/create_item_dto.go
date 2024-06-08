package item

import (
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ItemAvailability struct {
	Day  string `json:"day"`
	From string `json:"from"`
	To   string `json:"to"`
}
type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type LocalizationTextDesc struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}
type CreateItemDto struct {
	Id                string               `json:"_"`
	AccountId         string               `json:"account_id"`
	Name              LocalizationText     `json:"name" validate:"required,dive"`
	Desc              LocalizationTextDesc `json:"desc"`
	Type              string               `json:"type"`
	Min               int                  `json:"min"`
	Max               int                  `json:"max"`
	Calories          int                  `json:"calories" validate:"required"`
	Price             float64              `json:"price" validate:"required""`
	ModifierGroupsIds []string             `json:"modifier_groups_ids" validate:"modifier_groups_ids_rules"`
	Availabilities    []ItemAvailability   `json:"availabilities"`
	Tags              []string             `json:"tags"`
	Image             string               `json:"image" validate:"required"`
	Status            string               `json:"status" validate:"oneof=active inactive"`
	AdminDetails      []dto.AdminDetails   `json:"-"`
}

func (input *CreateItemDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validateModifierGroupsExistsInDB := func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().([]string)
		isValidObjectIds := utils.ValidateIDsIsMongoObjectIds(fl)
		if !isValidObjectIds {
			return false
		}
		return existsIModifierGroup("db", value)
	}
	validate.RegisterValidation("modifier_groups_ids_rules", validateModifierGroupsExistsInDB)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}

func existsIModifierGroup(db interface{}, value []string) bool {
	return true
}
