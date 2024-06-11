package custom_validators

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
)

type ItemCustomValidator struct {
	itemUsecase domain.ItemUseCase
}

func InitNewCustomValidatorsItem(itemUsecase domain.ItemUseCase) ItemCustomValidator {
	return ItemCustomValidator{
		itemUsecase: itemUsecase,
	}
}

func (i *ItemCustomValidator) ValidateNameIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		//read the value of the top struct to get accountId
		itemDto, ok := fl.Top().Interface().(*item.CreateItemDto)
		if !ok {
			fmt.Println("Unexpected type, expected *item.CreateItemDto")
			return false
		}
		accountId := itemDto.AccountId
		isExists, err := i.itemUsecase.CheckExists(context.Background(), accountId, val)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
