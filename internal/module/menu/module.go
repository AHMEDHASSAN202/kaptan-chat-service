package menu

import (
	"samm/internal/module/menu/delivery"
	"samm/internal/module/menu/repository/mongodb/item"
	menu_group2 "samm/internal/module/menu/repository/mongodb/menu_group"
	"samm/internal/module/menu/repository/mongodb/menu_group_item"
	"samm/internal/module/menu/repository/mongodb/modifier_group"
	useCaseItem "samm/internal/module/menu/usecase/item"
	"samm/internal/module/menu/usecase/menu_group"
	useCaseModifierGroup "samm/internal/module/menu/usecase/modifier_group"

	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		item.NewItemRepository,
		useCaseItem.NewItemUseCase,
		useCaseModifierGroup.NewModifierGroupUseCase,
		modifier_group.NewModifierGroupRepository,
		menu_group_item.NewMenuGroupItemRepository,
		menu_group2.NewMenuGroupRepository,
		menu_group.NewMenuGroupUseCase,
	),
	fx.Invoke(
		delivery.InitMenuGroupController,
		delivery.InitItemController,
		delivery.InitModifierGroupController,
	),
)
