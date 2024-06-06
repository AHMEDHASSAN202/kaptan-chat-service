package example

import (
	"go.uber.org/fx"
	"samm/internal/module/example/delivery"
	"samm/internal/module/example/repository/mongodb"
	"samm/internal/module/example/usecase"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewMenuMongoRepository,
		usecase.NewMenuUseCase,
	),
	fx.Invoke(delivery.InitUserController),
)
