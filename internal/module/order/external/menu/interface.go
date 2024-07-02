package menu

import (
	"context"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external/menu/responses"
	"samm/pkg/validators"
)

type Interface interface {
	GetMenuItemsDetails(ctx context.Context, menuItems []order.MenuItem) (responses.MenuDetailsResponse, validators.ErrorResponse)
}
