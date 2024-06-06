package location

import (
	"go.uber.org/fx"
	"samm/internal/module/location/delivery"
	"samm/internal/module/location/repository/location/mongodb"
	"samm/internal/module/location/usecase/location"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewLocationMongoRepository,
		location.NewLocationUseCase,
	),
	fx.Invoke(delivery.InitController),
)
