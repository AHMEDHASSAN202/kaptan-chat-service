package mobile

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/repository/structs/menu_group"
	"samm/internal/module/menu/responses/menu_group/mobile"
)

func GetMenuGroupItemBuilder(item *menu_group.MobileGetItem) *mobile.MobileGetItemResponse {
	itemResponse := mobile.MobileGetItemResponse{}
	itemResponse.ID = item.ID
	itemResponse.ItemId = item.ItemId
	itemResponse.Name.Ar = item.Name.Ar
	itemResponse.Name.En = item.Name.En
	itemResponse.Desc.Ar = item.Desc.Ar
	itemResponse.Desc.En = item.Desc.En
	itemResponse.Calories = item.Calories
	itemResponse.Price = item.Price
	itemResponse.ModifierGroupIds = item.ModifierGroupIds
	itemResponse.Tags = item.Tags
	itemResponse.Image = item.Image
	itemResponse.Category.ID = item.Category.ID
	itemResponse.Category.Name.Ar = item.Category.Name.Ar
	itemResponse.Category.Name.En = item.Category.Name.En
	itemResponse.Category.Icon = item.Category.Icon
	itemResponse.ModifierGroups = make([]mobile.ModifierGroupResponse, 0)
	addons := addonsAsMap(item.Addons)
	if item.ModifierGroups != nil {
		for _, modifierGroup := range item.ModifierGroups {
			group := mobile.ModifierGroupResponse{}
			group.ID = modifierGroup.ID
			group.Type = modifierGroup.Type
			group.Name.Ar = modifierGroup.Name.Ar
			group.Name.En = modifierGroup.Name.En
			group.Min = modifierGroup.Min
			group.Max = modifierGroup.Max
			group.ProductIds = modifierGroup.ProductIds
			group.Addons = make([]mobile.MobileGetItemAddonResponse, 0)
			if modifierGroup.ProductIds != nil {
				for _, addonId := range modifierGroup.ProductIds {
					if addon, ok := addons[addonId]; ok {
						group.Addons = append(group.Addons, addon)
					}
				}
			}
			itemResponse.ModifierGroups = append(itemResponse.ModifierGroups, group)
		}
	}
	return &itemResponse
}

func addonsAsMap(addons []menu_group.MobileGetItemAddon) map[primitive.ObjectID]mobile.MobileGetItemAddonResponse {
	addonsMap := map[primitive.ObjectID]mobile.MobileGetItemAddonResponse{}
	if addons != nil {
		for _, addon := range addons {
			addonsMap[addon.ID] = mobile.MobileGetItemAddonResponse{
				ID:       addon.ID,
				Name:     mobile.LocalizationText{Ar: addon.Name.Ar, En: addon.Name.En},
				Type:     addon.Type,
				Min:      addon.Min,
				Max:      addon.Max,
				Calories: addon.Calories,
				Price:    addon.Price,
				Image:    addon.Image,
			}
		}
	}
	return addonsMap
}
