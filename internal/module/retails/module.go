package retails

import (
	"go.uber.org/fx"
	"samm/internal/module/retails/delivery"
	account_repo "samm/internal/module/retails/repository/account/mongodb"
	brand_repo "samm/internal/module/retails/repository/brand"
	cuisine_repo "samm/internal/module/retails/repository/cuisine"
	"samm/internal/module/retails/repository/location/mongodb"
	"samm/internal/module/retails/usecase/account"
	brand_usecase "samm/internal/module/retails/usecase/brand"
	cuisine_usecase "samm/internal/module/retails/usecase/cuisine"
	"samm/internal/module/retails/usecase/location"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewLocationMongoRepository,
		location.NewLocationUseCase,
		cuisine_repo.NewCuisineRepository,
		cuisine_usecase.NewCuisineUseCase,
		brand_repo.NewBrandRepository,
		brand_usecase.NewBrandUseCase,
		account_repo.NewAccountMongoRepository,
		account.NewAccountUseCase,
	),
	fx.Invoke(delivery.InitController, delivery.InitCuisineController, delivery.InitBrandController, delivery.InitAccountController),
)
