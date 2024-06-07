package retails

import (
	"go.uber.org/fx"
	"samm/internal/module/retails/delivery"
	"samm/internal/module/retails/repository/location/mongodb"
	"samm/internal/module/retails/usecase/location"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewLocationMongoRepository,
		location.NewLocationUseCase,
	),
	fx.Invoke(delivery.InitController),
)
