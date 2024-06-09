package item

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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
	Name              LocalizationText     `json:"name" validate:"required"`
	Desc              LocalizationTextDesc `json:"desc"`
	Type              string               `json:"type"`
	Min               int                  `json:"min"`
	Max               int                  `json:"max"`
	Calories          int                  `json:"calories" validate:"required"`
	Price             float64              `json:"price" validate:"required"`
	ModifierGroupsIds []string             `json:"modifier_groups_ids" validate:"Modifier_groups_ids_rules"`
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

	// Register custom field-specific messages
	ctx := c.Request().Context()
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Modifier_groups_ids_rules,
		RegisterValidationFunc: validateModifierGroupsExistsInDB,
	})
}

func existsIModifierGroup(db interface{}, value []string) bool {
	return true
}
