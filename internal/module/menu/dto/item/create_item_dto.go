package item

import (
	"context"
	"github.com/go-playground/validator/v10"
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
	Ar string `json:"ar" validate:"required,min=3,Item_name_is_unique_rules_validation"`
	En string `json:"en" validate:"required,min=3,Item_name_is_unique_rules_validation"`
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
	ModifierGroupsIds []string             `json:"modifier_groups_ids" validate:"Invalid_mongo_ids_validation_rule"`
	Availabilities    []ItemAvailability   `json:"availabilities"`
	Tags              []string             `json:"tags"`
	Image             string               `json:"image" validate:"required"`
	Status            string               `json:"status" validate:"oneof=active inactive"`
	AdminDetails      []dto.AdminDetails   `json:"-"`
}

func (input *CreateItemDto) Validate(ctx context.Context, validate *validator.Validate, validateNameIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Item_name_is_unique_rules_validation,
		RegisterValidationFunc: validateNameIsUnique,
	})
}
