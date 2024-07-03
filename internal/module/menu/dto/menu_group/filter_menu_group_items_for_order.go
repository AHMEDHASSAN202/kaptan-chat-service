package menu_group

import (
	"samm/pkg/utils/dto"
)

type MenuItem struct {
	Id          string     `json:"id" validate:"required"`
	Qty         string     `json:"qty" validate:"required,min=1"`
	ModifierIds []MenuItem `json:"modifier_ids"`
}

type FilterMenuGroupItemsForOrder struct {
	MenuItems  []MenuItem `json:"menu_items" validate:"required,dive,min=1"`
	LocationId string     `json:"branch_id" validate:"required,mongodb"`
	dto.MobileHeaders
}
