package custom_validators

import (
	"context"
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
		var accountId, itemId string

		IitemDto := fl.Top().Interface()
		switch fl.Top().Type().String() {
		case "*item.CreateItemDto":
			itemDto := IitemDto.(*item.CreateItemDto)
			accountId = itemDto.AccountId
			itemId = itemDto.Id

		case "*item.CreateBulkItemDto":
			itemDto := IitemDto.(*item.CreateBulkItemDto)
			accountId = itemDto.AccountId
			itemId = itemDto.Id

		case "*item.UpdateItemDto":
			itemDto := IitemDto.(*item.UpdateItemDto)
			accountId = itemDto.AccountId
			itemId = itemDto.Id

		default:
			return false
		}

		isExists, err := i.itemUsecase.CheckExists(context.Background(), accountId, val, itemId)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
