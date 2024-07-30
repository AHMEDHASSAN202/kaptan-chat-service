package custom_validators

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/retails/domain"
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

func (i *RetailCustomValidator) ValidateAccountIsExists() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		accountId := fl.Field().Interface().(string)
		fmt.Println("accountId=> ", accountId)
		isExists, _ := i.accountUseCase.CheckAccountExists(context.Background(), accountId)
		fmt.Println("isExists=> ", isExists)
		return isExists
	}
}

func (i *RetailCustomValidator) ValidateCuisineIdsExists(cuisineIds []string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		//val := fl.Field().Interface().(string)
		//brandDto, ok := fl.Top().Interface().(*brand.CreateBrandDto)
		//if !ok {
		//	i.logger.Error("Unexpected type, expected *item.CreateBrandDto")
		//	return false
		//}
		fmt.Println("cuisineIds=> ", cuisineIds)
		ctx := context.Background()
		err := i.cuisineUsecase.CheckExists(&ctx, cuisineIds)
		if err.IsError {
			i.logger.Error(err.ErrorMessageObject)
			return false
		}
		return true
	}
}

func (i *RetailCustomValidator) ValidateCuisineNameUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		cuisineName := fl.Field().Interface().(string)
		isExists, err := i.cuisineUsecase.CheckNameExists(context.Background(), cuisineName)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
