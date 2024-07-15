package custom_validators

import (
	"github.com/go-playground/validator/v10"
	"samm/internal/module/kitchen/dto/kitchen"
)

type KitchenCustomValidator struct {
}

func InitNewCustomValidatorsKitchen() KitchenCustomValidator {
	return KitchenCustomValidator{}
}

func (i *KitchenCustomValidator) ValidateAccountAndLocationRequired() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {

		switch fl.Top().Interface().(type) {
		case *kitchen.StoreKitchenDto:
			dto, ok := fl.Top().Interface().(*kitchen.StoreKitchenDto)
			if !ok {
				return false
			}
			if len(dto.LocationIds) == 0 && len(dto.AccountIds) == 0 {
				return false
			}
			return true
		case *kitchen.UpdateKitchenDto:
			dto, ok := fl.Top().Interface().(*kitchen.UpdateKitchenDto)
			if !ok {
				return false
			}
			if len(dto.LocationIds) == 0 && len(dto.AccountIds) == 0 {
				return false
			}
			return true
		}
		return false
	}
}
