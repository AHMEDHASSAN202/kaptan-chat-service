package dashboard

import (
	"samm/internal/module/menu/repository/structs/menu_group"
	"samm/pkg/utils"
)

func FindMenuGroupBuilder(model *menu_group.FindMenuGroupWithItems) *menu_group.FindMenuGroupWithItems {
	if model == nil {
		return model
	}
	categories := make([]menu_group.MenuCategory, 0)
	for _, category := range model.Categories {
		if category.ID.IsZero() {
			continue
		}
		utils.SortByField(&category.Items, "Sort")
		categories = append(categories, category)
	}
	model.Categories = categories
	return model
}
