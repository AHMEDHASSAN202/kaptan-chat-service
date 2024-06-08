package retails

import (
	"go.uber.org/fx"
	"samm/internal/module/retails/delivery"
	cuisine_repo "samm/internal/module/retails/repository/cuisine"
	"samm/internal/module/retails/repository/location/mongodb"
	"samm/internal/module/retails/usecase/cuisine"
	"samm/internal/module/retails/usecase/location"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewLocationMongoRepository,
		location.NewLocationUseCase,
		cuisine_repo.NewCuisineRepository,
		cuisine.NewCuisineUseCase,
	),
	fx.Invoke(delivery.InitController, delivery.InitCuisineController),
)
