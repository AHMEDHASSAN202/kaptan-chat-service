package user

import (
	"go.uber.org/fx"
	"samm/internal/module/user/delivery"
	collection_method_repo "samm/internal/module/user/repository/collection_method"
	user_repo "samm/internal/module/user/repository/user"
	"samm/internal/module/user/usecase/collection_method"
	user_usecase "samm/internal/module/user/usecase/user"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		user_repo.NewUserMongoRepository,
		user_usecase.NewUserUseCase,
		collection_method.NewCollectionMethodUseCase,
		collection_method_repo.NewCollectionMethodMongoRepository,
	),
	fx.Invoke(
		delivery.InitUserController,
		delivery.InitCollectionMethodController,
	),
)
