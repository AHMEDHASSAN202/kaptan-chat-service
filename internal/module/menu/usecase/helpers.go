package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/utils"
)

func (oRec *MenuGroupUseCase) InjectItemsToDTO(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) error {
	itemIds := []primitive.ObjectID{}
	if dto.Categories != nil {
		for _, category := range dto.Categories {
			if category.MenuItems != nil {
				for _, menuItem := range category.MenuItems {
					itemIds = append(itemIds, utils.ConvertStringIdToObjectId(menuItem.ItemId))
				}
			}
		}
	}

	items, err := oRec.itemRepo.GetByIds(ctx, itemIds)
	if err != nil || items == nil {
		return err
	}

	itemsMap := map[string]domain.Item{}
	for _, item := range items {
		itemsMap[item.ID.String()] = item
	}

	if dto.Categories != nil {
		for i, category := range dto.Categories {
			if category.MenuItems != nil {
				for ii, menuItem := range category.MenuItems {
					item := itemsMap[menuItem.ItemId]
					menuGroupItem := menu_group.MenuItemDTO{}
					if menuItem.Id == "" || utils.IsValidateObjectId(menuItem.Id) {
						menuItem.Id = utils.ConvertObjectIdToStringId(primitive.NewObjectID())
					}
					menuGroupItem.ItemId = utils.ConvertObjectIdToStringId(item.ID)
					menuGroupItem.Name.En = item.Name.En
					menuGroupItem.Name.Ar = item.Name.Ar
					menuGroupItem.Desc.En = item.Desc.En
					menuGroupItem.Desc.Ar = item.Desc.Ar
					menuGroupItem.Calories = item.Calories
					menuGroupItem.Price = utils.If(menuItem.Price == 0, item.Price, menuItem.Price).(float64)
					menuGroupItem.ModifierGroupsIds = item.ModifierGroupsIds
					menuGroupItem.Tags = item.Tags
					menuGroupItem.Status = menuItem.Status
					menuGroupItem.Image = item.Image
					menuGroupItem.Sort = menuItem.Sort
					menuGroupItem.Availabilities = []menu_group.AvailabilityDTO{}
					if item.Availabilities != nil {
						for _, availability := range item.Availabilities {
							menuGroupItem.Availabilities = append(menuGroupItem.Availabilities, menu_group.AvailabilityDTO{
								Day:  availability.Day,
								From: availability.From,
								To:   availability.To,
							})
						}
					}
					dto.Categories[i].MenuItems[ii] = menuGroupItem
				}
			}
		}
	}

	return nil
}
