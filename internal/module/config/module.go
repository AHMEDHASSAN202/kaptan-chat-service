package config

import (
	"samm/internal/module/config/custom_validators"
	"samm/internal/module/config/delivery"
	"samm/internal/module/config/repository/app_config"

	appConfigUsecase "samm/internal/module/config/usecase/app_config"

	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		app_config.NewAppConfigRepository,
		appConfigUsecase.NewAppConfigUseCase,
		custom_validators.InitNewCustomValidatorsAppConfig,
	),
	fx.Invoke(
		delivery.InitAppConfigController,
	),
)
