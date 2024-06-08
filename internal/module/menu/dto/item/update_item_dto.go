package item

import (
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UpdateItemDto struct {
	Id                string               `json:"_"`
	AccountId         string               `json:"account_id" `
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

func (input *UpdateItemDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validate.RegisterValidation("modifier_groups_ids_rules", utils.ValidateIDsIsMongoObjectIds)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
