package mobile

import (
	"samm/internal/module/menu/repository/structs/menu_group"
	"samm/internal/module/menu/responses/menu_group/mobile"
)

func GetMenuGroupItemsBuilder(items *[]menu_group.MobileGetMenuGroupItems) *[]mobile.GetMenuGroupItemsResponse {
	result := make([]mobile.GetMenuGroupItemsResponse, 0)
	if items == nil {
		return &result
	}
	for _, groupItems := range *items {
		category := mobile.GetMenuGroupItemsResponse{
			ID:    groupItems.ID,
			Icon:  groupItems.Icon,
			Sort:  groupItems.Sort,
			Items: make([]mobile.GetMenuGroupItem, 0),
			Name:  mobile.LocalizationText{Ar: groupItems.Name.Ar, En: groupItems.Name.En},
		}
		if groupItems.Items != nil {
			for _, item := range groupItems.Items {
				category.Items = append(category.Items, mobile.GetMenuGroupItem{
					ID:       item.ID,
					Name:     mobile.LocalizationText{Ar: item.Name.Ar, En: item.Name.En},
					Image:    item.Image,
					Price:    item.Price,
					Calories: item.Calories,
					Sort:     item.Sort,
				})
			}
		}
		result = append(result, category)
	}
	return &result
}
