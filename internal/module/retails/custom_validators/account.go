package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/logger"
)

type RetailCustomValidator struct {
	accountUseCase domain.AccountUseCase
	cuisineUsecase domain.CuisineUseCase
	logger         logger.ILogger
}

func InitNewCustomValidators(accountUseCase domain.AccountUseCase, cuisineUsecase domain.CuisineUseCase, log logger.ILogger) RetailCustomValidator {
	return RetailCustomValidator{
		accountUseCase: accountUseCase,
		cuisineUsecase: cuisineUsecase,
		logger:         log,
	}
}

func (i *RetailCustomValidator) ValidateAccountEmailIsUnique(accountId string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().Interface().(string)
		isExists := i.accountUseCase.CheckAccountEmail(context.Background(), email, accountId)
		return !isExists
	}
}

func (i *RetailCustomValidator) ValidateCuisineIdsExists() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		//val := fl.Field().Interface().(string)
		brandDto, ok := fl.Top().Interface().(*brand.CreateBrandDto)
		if !ok {
			i.logger.Error("Unexpected type, expected *item.CreateBrandDto")
			return false
		}
		ctx := context.Background()
		err := i.cuisineUsecase.CheckExists(&ctx, brandDto.CuisineIds)
		if err.IsError {
			i.logger.Error(err.ErrorMessageObject)
			return false
		}
		return true
	}
}
