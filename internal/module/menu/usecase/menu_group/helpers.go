package menu_group

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (oRec *MenuGroupUseCase) InjectItemsToDTO(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) validators.ErrorResponse {
	itemIds := []primitive.ObjectID{}
	if dto.Categories != nil {
		for _, category := range dto.Categories {
			if category.MenuItems != nil {
				for _, menuItem := range category.MenuItems {
					itemId := utils.ConvertStringIdToObjectId(menuItem.ItemId)
					if !utils.Contains(itemIds, itemId) {
						itemIds = append(itemIds, itemId)
					}
				}
			}
		}
	}

	items, err := oRec.itemRepo.GetByIds(ctx, itemIds)
	spew.Dump("items", items)
	if err != nil || items == nil {
		oRec.logger.Error(err)
		return validators.GetErrorResponse(&ctx, "E1000", nil, nil)
	}

	itemsMap := map[string]domain.Item{}
	for _, item := range items {
		itemsMap[utils.ConvertObjectIdToStringId(item.ID)] = item
	}

	if len(itemsMap) != len(itemIds) {
		return validators.GetErrorResponse(&ctx, "E1001", nil, nil)
	}

	if dto.Categories != nil {
		for i, category := range dto.Categories {
			if category.MenuItems != nil {
				for ii, menuItem := range category.MenuItems {
					item := itemsMap[menuItem.ItemId]
					menuGroupItem := menu_group.MenuItemDTO{}
					if menuItem.Id == "" || !utils.IsValidateObjectId(menuItem.Id) {
						menuItem.Id = utils.ConvertObjectIdToStringId(primitive.NewObjectID())
						menuGroupItem.IsNew = true
					}
					menuGroupItem.Id = menuItem.Id
					menuGroupItem.ItemId = utils.ConvertObjectIdToStringId(item.ID)
					menuGroupItem.Name.En = item.Name.En
					menuGroupItem.Name.Ar = item.Name.Ar
					menuGroupItem.Desc.En = item.Desc.En
					menuGroupItem.Desc.Ar = item.Desc.Ar
					menuGroupItem.Calories = item.Calories
					menuGroupItem.Price = utils.If(menuItem.Price == 0, item.Price, menuItem.Price).(float64)
					menuGroupItem.ModifierGroupIds = item.ModifierGroupIds
					menuGroupItem.Tags = item.Tags
					menuGroupItem.Status = menuItem.Status
					menuGroupItem.Image = item.Image
					menuGroupItem.Sort = menuItem.Sort
					menuGroupItem.ApprovalData.HasOriginal = item.ApprovalData.HasOriginal
					menuGroupItem.ApprovalData.ApprovalStatus = item.ApprovalData.ApprovalStatus
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
	return validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) AuthorizeMenuGroup(ctx *context.Context, menuGroup *domain.MenuGroup, accountId primitive.ObjectID) validators.ErrorResponse {
	if menuGroup == nil || menuGroup.ID.IsZero() {
		oRec.logger.Error("AuthorizeMenuGroup-> Error In menuGroup")
		return validators.GetErrorResponse(ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	if !menuGroup.Authorized(accountId) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Menu Group -> ", menuGroup.ID)
		return validators.GetErrorResponse(ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}
	return validators.ErrorResponse{}
}
