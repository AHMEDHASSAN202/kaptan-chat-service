package menu

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/menu/dto/menu_group"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external/menu/responses"
	"samm/pkg/validators"
)

func (i IService) GetMenuItemsDetails(ctx context.Context, menuItems []order.MenuItem, locationId string) ([]responses.MenuDetailsResponse, validators.ErrorResponse) {

	input := menu_group.FilterMenuGroupItemsForOrder{}
	copier.Copy(&input.MenuItems, menuItems)
	input.LocationId = locationId
	menus, err := i.MenuUseCase.MobileFilterMenuGroupItemForOrder(ctx, &input)
	if err != nil {
		return []responses.MenuDetailsResponse{}, validators.GetErrorResponseFromErr(err)
	}
	resp := make([]responses.MenuDetailsResponse, 0)
	copier.Copy(&resp, &menus)

	return resp, validators.ErrorResponse{}
}
