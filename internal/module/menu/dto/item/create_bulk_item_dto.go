package item

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type CreateBulkItemDto struct {
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
	Image             string               `json:"image"`
	Status            string               `json:"status" validate:"oneof=active inactive"`
	AdminDetails      []dto.AdminDetails   `json:"-"`
}

func (input *CreateBulkItemDto) validateItem(ctx context.Context, rowIndex int, validate *validator.Validate, validateNameIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	// Register custom field-specific messages
	validationErr := validators.ValidateStruct(ctx, validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Item_name_is_unique_rules_validation,
		RegisterValidationFunc: validateNameIsUnique,
	}, validators.CustomErrorTags{
		ValidationTag:          localization.Modifier_items_cant_contains_modifier_group,
		RegisterValidationFunc: ModifierItemsCantContainsModifierGroup,
	})
	//prepare the validation response to be array
	if validationErr.IsError {
		for k, value := range validationErr.ValidationErrors {
			delete(validationErr.ValidationErrors, k)
			validationErr.ValidationErrors[fmt.Sprintf("product.%d.%s", rowIndex, k)] = value
		}
	}
	return validationErr
}
func (input *CreateBulkItemDto) Validate(ctx context.Context, payload []CreateBulkItemDto, validate *validator.Validate, validateNameIsUnique func(fl validator.FieldLevel) bool) validators.ErrorResponse {

	validationErrs := validators.ErrorResponse{ValidationErrors: map[string][]string{}}

	if len(payload) == 0 {
		return validators.GetErrorResponse(&ctx, localization.E1001, nil, nil)
	}
	for index, itemDoc := range payload {
		validationErr := itemDoc.validateItem(ctx, index, validate, validateNameIsUnique)
		if validationErr.IsError {
			validationErrs.IsError = validationErr.IsError
			for k, v := range validationErr.ValidationErrors {
				validationErrs.ValidationErrors[k] = v
			}
		}
	}

	return validationErrs
}
