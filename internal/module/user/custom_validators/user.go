package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/user/domain"
	"samm/pkg/logger"
)

type UserCustomValidator struct {
	userUseCase domain.UserUseCase
	logger      logger.ILogger
}

func InitNewCustomValidatorsForUser(userUseCase domain.UserUseCase, log logger.ILogger) UserCustomValidator {
	return UserCustomValidator{
		userUseCase: userUseCase,
		logger:      log,
	}
}

func (i *UserCustomValidator) ValidateUserEmailIsUnique(userId string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().Interface().(string)
		isExists := i.userUseCase.UserEmailExists(context.Background(), email, userId)
		return !isExists
	}
}
