package menu

import (
	"samm/internal/module/menu/delivery"
	"samm/internal/module/menu/repository/mongodb"
	"samm/internal/module/menu/repository/mongodb/item"
	"samm/internal/module/menu/repository/mongodb/modifier_group"
	"samm/internal/module/menu/repository/mongodb/sku"
	"samm/internal/module/menu/usecase"
	useCaseItem "samm/internal/module/menu/usecase/item"
	useCaseModifierGroup "samm/internal/module/menu/usecase/modifier_group"
	useCaseSku "samm/internal/module/menu/usecase/sku"

	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		item.NewItemRepository,
		useCaseItem.NewItemUseCase,
		useCaseModifierGroup.NewModifierGroupUseCase,
		modifier_group.NewModifierGroupRepository,
		mongodb.NewMenuGroupItemRepository,
		mongodb.NewMenuGroupRepository,
		usecase.NewMenuGroupUseCase,
		useCaseSku.NewSKUUseCase,
		sku.NewSkuRepository,
	),
	fx.Invoke(
		delivery.InitMenuGroupController,
		delivery.InitItemController,
		delivery.InitModifierGroupController,
		delivery.InitSKUController,
	),
)
