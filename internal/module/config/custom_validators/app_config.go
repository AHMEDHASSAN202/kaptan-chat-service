package custom_validators

import (
	"context"
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"

	"github.com/go-playground/validator/v10"
)

type AppConfigCustomValidator struct {
	appConfigUsecase domain.AppConfigUseCase
}

func InitNewCustomValidatorsAppConfig(appConfigUsecase domain.AppConfigUseCase) AppConfigCustomValidator {
	return AppConfigCustomValidator{
		appConfigUsecase: appConfigUsecase,
	}
}

func (i *AppConfigCustomValidator) ValidateAppTypeIsUnique() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		appConfigDto := fl.Top().Interface().(*app_config.CreateUpdateAppConfigDto)
		isExists, err := i.appConfigUsecase.CheckExists(context.Background(), val, appConfigDto.Id)
		if err.IsError {
			return false
		}
		return !isExists
	}
}
