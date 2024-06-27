package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
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

func (i *UserCustomValidator) ValidateUserEmailIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().Interface().(string)
		profile, ok := fl.Top().Interface().(*user.UpdateUserProfileDto)
		if !ok {
			return false
		}
		userId := profile.ID
		ctx := context.Background()
		isExists := i.userUseCase.UserEmailExists(&ctx, email, userId)
		return !isExists
	}
}
