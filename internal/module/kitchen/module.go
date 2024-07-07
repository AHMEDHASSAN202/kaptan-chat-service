package kitchen

import (
	"go.uber.org/fx"
	"samm/internal/module/kitchen/delivery"
	kitchen_repo "samm/internal/module/kitchen/repository/kitchen"
	kitchen_usecase "samm/internal/module/kitchen/usecase/kitchen"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		kitchen_repo.NewKitchenMongoRepository,
		kitchen_usecase.NewKitchenUseCase,
	),
	fx.Invoke(
		delivery.InitKitchenController,
	),
)


