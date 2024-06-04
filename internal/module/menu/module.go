package menu

import (
	"example.com/fxdemo/internal/module/menu/delivery"
	"example.com/fxdemo/internal/module/menu/repository/mongodb"
	"example.com/fxdemo/internal/module/menu/usecase"
	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewMenuMongoRepository,
		usecase.NewMenuUseCase,
	),
	fx.Invoke(delivery.InitUserController),
)
