package menu

import (
	"go.uber.org/fx"
	"samm/internal/module/menu/delivery"
	"samm/internal/module/menu/repository/mongodb"
	"samm/internal/module/menu/repository/mongodb/item"
	"samm/internal/module/menu/usecase"
	useCaseItem "samm/internal/module/menu/usecase/item"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		item.NewItemRepository,
		useCaseItem.NewItemUseCase,
		mongodb.NewMenuGroupItemRepository,
		mongodb.NewMenuGroupRepository,
		usecase.NewMenuGroupUseCase,
	),
	fx.Invoke(
		delivery.InitMenuGroupController,
		delivery.InitItemController,
	),
)
