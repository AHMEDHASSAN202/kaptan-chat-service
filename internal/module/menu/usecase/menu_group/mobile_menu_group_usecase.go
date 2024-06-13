package menu_group

import (
	"context"
	"net/http"
	"samm/internal/module/menu/builder/menu_group/mobile"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (oRec *MenuGroupUseCase) MobileGetMenuGroupItems(ctx context.Context, dto menu_group.GetMenuGroupItemsDTO) (interface{}, validators.ErrorResponse) {
	items, err := oRec.menuGroupItemRepo.MobileGetMenuGroupItems(ctx, &dto)
	result := mobile.GetMenuGroupItemsBuilder(items)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> MobileGetMenuGroupItems -> ", err)
		return result, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return result, validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) MobileGetMenuGroupItem(ctx context.Context, dto menu_group.GetMenuGroupItemDTO) (interface{}, validators.ErrorResponse) {
	item, err := oRec.menuGroupItemRepo.MobileGetMenuGroupItem(ctx, &dto)
	result := mobile.GetMenuGroupItemBuilder(item)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> MobileGetMenuGroupItem -> ", err)
		return result, validators.GetErrorResponse(&ctx, localization.E1002Item, nil, utils.GetAsPointer(http.StatusNotFound))
	}
	return result, validators.ErrorResponse{}
}
