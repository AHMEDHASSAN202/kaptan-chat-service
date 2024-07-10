package item

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type Availability struct {
	Day  string `json:"day" validate:"required,oneof=mon tue wed thu fri sat sun"`
	From string `json:"from" validate:"required,Timeformat"`
	To   string `json:"to" validate:"required,Timeformat"`
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
	Type              string               `json:"type" validate:"required,oneof=product modifier"`
	Min               int                  `json:"min"`
	Max               int                  `json:"max"`
	Calories          int                  `json:"calories" validate:"required"`
	Price             float64              `json:"price" validate:"required"`
	ModifierGroupsIds []string             `json:"modifier_groups_ids" validate:"Invalid_mongo_ids_validation_rule,Modifier_items_cant_contains_modifier_group"`
	Availabilities    []Availability       `json:"availabilities"`
	Tags              []string             `json:"tags"`
	SKU               string               `json:"sku"`
	Image             string               `json:"image" validate:"required"`
	Status            string               `json:"status" validate:"oneof=active inactive"`
	AdminDetails      []dto.AdminDetails   `json:"-"`
}

var ModifierItemsCantContainsModifierGroup = func(fl validator.FieldLevel) bool {
	IitemDto := fl.Top().Interface()
	val := fl.Field().Interface().([]string)
	//read the value of the top struct to get accountId
	var itemType string

	switch fl.Top().Type().String() {
	case "*item.CreateItemDto":
		itemDto := IitemDto.(*CreateItemDto)
		itemType = itemDto.Type

	case "*item.CreateBulkItemDto":
		itemDto := IitemDto.(*CreateBulkItemDto)
		itemType = itemDto.Type

	case "*item.UpdateItemDto":
		itemDto := IitemDto.(*UpdateItemDto)
		itemType = itemDto.Type

	default:
		return false
	}
	if itemType == "modifier" && len(val) > 0 {
		return false
	} else {
		return true
	}
}

func (input *CreateItemDto) Validate(ctx context.Context, validate *validator.Validate, validateNameIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	// Register custom field-specific messages
	return validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Item_name_is_unique_rules_validation,
		RegisterValidationFunc: validateNameIsUnique,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Modifier_items_cant_contains_modifier_group,
		RegisterValidationFunc: ModifierItemsCantContainsModifierGroup,
	})
}
