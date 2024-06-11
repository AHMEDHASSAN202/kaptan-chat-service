package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/retails/domain"
)

type RetailCustomValidator struct {
	accountUseCase domain.AccountUseCase
}

func InitNewCustomValidators(accountUseCase domain.AccountUseCase) RetailCustomValidator {
	return RetailCustomValidator{
		accountUseCase: accountUseCase,
	}
}

func (i *RetailCustomValidator) ValidateAccountEmailIsUnique(accountId string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().Interface().(string)
		isExists := i.accountUseCase.CheckAccountEmail(context.Background(), email, accountId)
		return !isExists
	}
}
