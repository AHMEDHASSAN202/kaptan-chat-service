package custom_validators

import (
	"context"
	"samm/internal/module/menu/domain"

	"github.com/go-playground/validator/v10"
)

type SKUCustomValidator struct {
	skuUsecase domain.SKUUseCase
}

func InitNewCustomValidatorsSKU(skuUsecase domain.SKUUseCase) SKUCustomValidator {
	return SKUCustomValidator{
		skuUsecase: skuUsecase,
	}
}

func (i *SKUCustomValidator) ValidateSKUIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		isExists, err := i.skuUsecase.CheckExists(context.Background(), val)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
