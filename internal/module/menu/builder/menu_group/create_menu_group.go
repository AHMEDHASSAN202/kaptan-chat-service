package menu_group

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/utils"
)

func MenuGroupBuilder(dto *menu_group.CreateMenuGroupDTO) (*domain.MenuGroup, *[]domain.MenuGroupItem) {
	menuGroupDomain := domain.MenuGroup{}
	if dto.ID == primitive.NilObjectID || dto.ID.IsZero() {
		dto.ID = primitive.NewObjectID()
	}
	menuGroupDomain.AccountId = ""
	menuGroupDomain.Name.Ar = dto.Name.Ar
	menuGroupDomain.Name.En = dto.Name.En
	menuGroupDomain.BranchIds = utils.ConvertStringIdsToObjectIds(utils.RemoveDuplicates[string](dto.BranchIds))
	menuGroupDomain.Categories = CategoriesBuilder(&dto.Categories)
	menuGroupDomain.Availabilities = AvailabilitiesBuilder(dto.Availabilities)
	menuGroupDomain.Status = dto.Status
	items := MenuGroupItemsBuilder(dto)
	return &menuGroupDomain, items
}

func CategoriesBuilder(categoriesInput *[]menu_group.CategoryDTO) []domain.Category {
	categories := make([]domain.Category, 0)
	if categoriesInput != nil && len(*categoriesInput) >= 1 {
		for _, category := range *categoriesInput {
			cat := domain.Category{}
			if category.ID == "" || utils.IsValidateObjectId(category.ID) {
				category.ID = utils.ConvertObjectIdToStringId(primitive.NewObjectID())
			}
			cat.ID = utils.ConvertStringIdToObjectId(category.ID)
			cat.Name.En = category.Name.En
			cat.Name.Ar = category.Name.Ar
			cat.Icon = category.Icon
			cat.Sort = category.Sort
			cat.Status = category.Status
			cat.MenuItemIds = []primitive.ObjectID{}
			if category.MenuItems != nil {
				for _, item := range category.MenuItems {
					if item.Id == "" {
						continue
					}
					menuGroupItemId := utils.ConvertStringIdToObjectId(item.Id)
					if !utils.Contains(cat.MenuItemIds, menuGroupItemId) {
						cat.MenuItemIds = append(cat.MenuItemIds, menuGroupItemId)
					}
				}
			}
			categories = append(categories, cat)
		}
	}
	return categories
}

func AvailabilitiesBuilder(availabilitiesInput []menu_group.AvailabilityDTO) []domain.MenuGroupAvailability {
	availabilities := make([]domain.MenuGroupAvailability, 0)
	if availabilitiesInput != nil && len(availabilitiesInput) >= 1 {
		for _, availability := range availabilitiesInput {
			availabilities = append(availabilities, domain.MenuGroupAvailability{
				Day:  availability.Day,
				From: availability.From,
				To:   availability.To,
			})
		}
	}
	return availabilities
}

func MenuGroupItemsBuilder(dto *menu_group.CreateMenuGroupDTO) *[]domain.MenuGroupItem {
	items := []domain.MenuGroupItem{}
	menuGroup := domain.ItemMenuGroup{
		ID:             dto.ID,
		BranchIds:      utils.ConvertStringIdsToObjectIds(utils.RemoveDuplicates[string](dto.BranchIds)),
		Status:         dto.Status,
		Availabilities: AvailabilitiesBuilder(dto.Availabilities),
	}
	if dto.Categories != nil {
		for _, category := range dto.Categories {
			if category.MenuItems != nil {
				for _, item := range category.MenuItems {
					menuGroupItem := domain.MenuGroupItem{}
					menuGroupItem.ID = utils.ConvertStringIdToObjectId(item.Id)
					menuGroupItem.ItemId = utils.ConvertStringIdToObjectId(item.ItemId)
					menuGroupItem.Name.En = item.Name.En
					menuGroupItem.Name.Ar = item.Name.Ar
					menuGroupItem.Desc.En = item.Desc.En
					menuGroupItem.Desc.Ar = item.Desc.Ar
					menuGroupItem.Calories = item.Calories
					menuGroupItem.Price = item.Price
					menuGroupItem.ModifierGroupsIds = item.ModifierGroupsIds
					menuGroupItem.Tags = item.Tags
					menuGroupItem.Status = item.Status
					menuGroupItem.Image = item.Image
					menuGroupItem.Category = domain.MenuGroupItemCategory{
						ID: utils.ConvertStringIdToObjectId(category.ID),
						Name: domain.LocalizationText{
							En: category.Name.En,
							Ar: category.Name.Ar,
						},
						Status: category.Status,
						Sort:   category.Sort,
						Icon:   category.Icon,
					}
					menuGroupItem.Availabilities = []domain.ItemAvailability{}
					if item.Availabilities != nil {
						for _, availability := range item.Availabilities {
							menuGroupItem.Availabilities = append(menuGroupItem.Availabilities, domain.ItemAvailability{
								Day:  availability.Day,
								From: availability.From,
								To:   availability.To,
							})
						}
					}
					menuGroupItem.MenuGroup = menuGroup
					items = append(items, menuGroupItem)
				}
			}
		}
	}
	return &items
}
