package helper

import (
	"context"
	"fmt"
	"net/http"
	"samm/internal/module/order/dto/order"
	extMenuResponses "samm/internal/module/order/external/menu/responses"
	extLocResponses "samm/internal/module/order/external/retails/responses"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func CheckIsLocationReadyForNewOrder(ctx *context.Context, doc extLocResponses.LocationDetails) validators.ErrorResponse {
	if !doc.IsOpen {
		return validators.GetErrorResponseWithErrors(ctx, localization.Mobile_location_not_open_error, nil)
	}
	return validators.ErrorResponse{}
}

func CheckIsMenuItemsValid(ctx *context.Context, menuDocs []extMenuResponses.MenuDetailsResponse, menuItemDto []order.MenuItem, createOrder bool) validators.ErrorResponse {
	modifierMap := make(map[string]extMenuResponses.MobileGetItemAddon)
	menuItemsMap := make(map[string]extMenuResponses.MenuDetailsResponse)
	qtyModifiers := make(map[string]map[string]int)

	for _, menu := range menuDocs {
		menuItemsMap[menu.ID.Hex()] = menu
		for _, addon := range menu.Addons {
			modifierMap[addon.ID.Hex()] = addon
			qtyModifiers[menu.ID.Hex()+"-"+addon.ID.Hex()] = map[string]int{"min": addon.Min, "max": addon.Max}
		}
	}

	ValidationErrors := make(map[string][]string)
	for menuIndex, item := range menuItemDto {
		if _, ok := menuItemsMap[item.Id]; !ok {
			ValidationErrors[fmt.Sprintf("menu.%d", menuIndex)] = []string{localization.GetTranslation(ctx, localization.Mobile_item_unavailable, nil, "")}
			continue
		}
		for modifierIndex, modifier := range item.ModifierIds {
			if _, ok := modifierMap[modifier.Id]; !ok {
				ValidationErrors[fmt.Sprintf("menu.%d.modifier.%d", menuIndex, modifierIndex)] = []string{localization.GetTranslation(ctx, localization.Mobile_item_unavailable, nil, "")}
				continue
			}
			if q, _ := qtyModifiers[item.Id+"-"+modifier.Id]; q["min"] != 0 && q["min"] > int(modifier.Qty) {
				ValidationErrors[fmt.Sprintf("menu.%d.modifier.%d", menuIndex, modifierIndex)] = []string{localization.GetTranslation(ctx, localization.MinModifierQty, map[string]interface{}{"modifier_name": localization.GetAttrByLang(ctx, modifierMap[modifier.Id].Name.En, modifierMap[modifier.Id].Name.Ar)}, "")}
			}
			if q, _ := qtyModifiers[item.Id+"-"+modifier.Id]; q["max"] != 0 && q["max"] < int(modifier.Qty) {
				ValidationErrors[fmt.Sprintf("menu.%d.modifier.%d", menuIndex, modifierIndex)] = []string{localization.GetTranslation(ctx, localization.MaxModifierQty, map[string]interface{}{"modifier_name": localization.GetAttrByLang(ctx, modifierMap[modifier.Id].Name.En, modifierMap[modifier.Id].Name.Ar)}, "")}
			}
		}
		validateModifierGroupQty(ctx, menuIndex, ValidationErrors, menuItemsMap[item.Id], item)
	}

	if len(ValidationErrors) >= 1 {
		return validators.ErrorResponse{IsError: true, ValidationErrors: ValidationErrors, StatusCode: http.StatusUnprocessableEntity}
	}

	return validators.ErrorResponse{}
}

func validateModifierGroupQty(ctx *context.Context, menuIndex int, validationErrors map[string][]string, menuDoc extMenuResponses.MenuDetailsResponse, menuItemDto order.MenuItem) {
	for _, group := range menuDoc.ModifierGroups {
		dtoGroupQty := 0
		for _, modifier := range menuItemDto.ModifierIds {
			if utils.Contains(group.ProductIds, utils.ConvertStringIdToObjectId(modifier.Id)) {
				dtoGroupQty += int(modifier.Qty)
			}
		}
		if group.Min != 0 && dtoGroupQty < group.Min {
			validationErrors[fmt.Sprintf("menu.%d.modifier_group.%d", menuIndex, menuIndex)] = []string{localization.GetTranslation(ctx, localization.MinModifierGroupQty, map[string]interface{}{"modifier_group_name": localization.GetAttrByLang(ctx, group.Name.En, group.Name.Ar)}, "")}
		}
		if group.Max != 0 && dtoGroupQty > group.Max {
			validationErrors[fmt.Sprintf("menu.%d.modifier_group.%d", menuIndex, menuIndex)] = []string{localization.GetTranslation(ctx, localization.MaxModifierGroupQty, map[string]interface{}{"modifier_group_name": localization.GetAttrByLang(ctx, group.Name.En, group.Name.Ar)}, "")}
		}
	}
}
