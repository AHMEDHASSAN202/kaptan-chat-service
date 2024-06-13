package menu_group

import (
	"context"
	"samm/internal/module/menu/builder/menu_group/mobile"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (oRec *MenuGroupUseCase) MobileGetMenuGroupItems(ctx context.Context, dto menu_group.GetMenuGroupItemDTO) (interface{}, validators.ErrorResponse) {
	items, err := oRec.menuGroupItemRepo.MobileGetMenuGroupItems(ctx, &dto)
	result := mobile.GetMenuGroupItemsBuilder(items)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> MobileGetMenuGroupItems -> ", err)
		return result, validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	return result, validators.ErrorResponse{}
}
