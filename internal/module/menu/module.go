package menu

import (
	"go.uber.org/fx"
	"samm/internal/module/menu/delivery"
	"samm/internal/module/menu/repository/mongodb"
	"samm/internal/module/menu/usecase"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mongodb.NewItemRepository,
		mongodb.NewMenuGroupRepository,
		usecase.NewMenuGroupUseCase,
	),
	fx.Invoke(
		delivery.InitMenuGroupController,
	),
)
