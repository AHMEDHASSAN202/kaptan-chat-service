package user

import (
	"go.uber.org/fx"
	"samm/internal/module/user/custom_validators"
	"samm/internal/module/user/delivery"
	user_repo "samm/internal/module/user/repository/user"
	user_usecase "samm/internal/module/user/usecase/user"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		user_repo.NewUserMongoRepository,
		user_usecase.NewUserUseCase,
		custom_validators.InitNewCustomValidatorsForUser,
	),
	fx.Invoke(
		delivery.InitUserController,
	),
)
