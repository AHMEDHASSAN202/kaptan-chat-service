package config

import (
	"samm/internal/module/config/delivery"
	"samm/internal/module/config/repository/app_config"

	appConfigUsecase "samm/internal/module/config/usecase/app_config"

	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		app_config.NewAppConfigRepository,
		appConfigUsecase.NewAppConfigUseCase,
	),
	fx.Invoke(
		delivery.InitAppConfigController,
	),
)
