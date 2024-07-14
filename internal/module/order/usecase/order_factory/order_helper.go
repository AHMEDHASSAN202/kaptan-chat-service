package order_factory

import (
	"context"
	"fmt"
	"net/http"
	"samm/internal/module/order/dto/order"
	extMenuResponses "samm/internal/module/order/external/menu/responses"
	extLocResponses "samm/internal/module/order/external/retails/responses"
	"samm/internal/module/order/responses"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func CalculateOrderCostBuilder(ctx context.Context, loc extLocResponses.LocationDetails, menus []extMenuResponses.MenuDetailsResponse, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	resp.Location = responses.LocationDoc{
		Id: loc.Id,
		Name: responses.LocalizationText{
			Ar: loc.Name.Ar,
			En: loc.Name.En,
		},
		IsOpen: loc.IsOpen,
	}
	var totalMenusValueBefore, totalMenusValueAfter float64
	resp.MenuItems, totalMenusValueBefore, totalMenusValueAfter = CalculateTotalCostForMenus(&ctx, menus, payload)
	resp.TotalPriceSummary = responses.TotalPriceSummary{
		Fees:                     0,
		TotalPriceBeforeDiscount: totalMenusValueBefore,
		TotalPriceAfterDiscount:  totalMenusValueAfter,
	}
	return resp, err
}

func CalculateTotalCostForMenus(ctx *context.Context, menus []extMenuResponses.MenuDetailsResponse, payload *order.CalculateOrderCostDto) (menuRespDocs []responses.MenuDoc, totalMenusValueBefore float64, totalMenusValueAfter float64) {
	//validate the menu&modifiers items
	validationErrorMap := CheckIsMenuItemsValid(ctx, menus, payload.MenuItems).ValidationErrors
	menuRespDocs = make([]responses.MenuDoc, 0)
	for _, item := range payload.MenuItems {
		for _, menu := range menus {
			if menu.ID.Hex() == item.Id {
				modifierDocs, totalModifierValueBefore, totalModifierValueAfter := GetModifierItemPriceSummary(item, menu, validationErrorMap)

				//calculate total values
				SubTotalMenusValueBefore := float64(item.Qty) * (menu.Price + totalModifierValueBefore)
				SubTotalMenusValueAfter := float64(item.Qty) * (menu.Price + totalModifierValueAfter)

				totalMenusValueBefore += SubTotalMenusValueBefore
				totalMenusValueAfter += SubTotalMenusValueAfter

				menuRespDoc := responses.MenuDoc{
					Id: item.Id,
					Name: responses.LocalizationText{
						Ar: menu.Name.Ar,
						En: menu.Name.En,
					},
					Desc: responses.LocalizationText{
						Ar: menu.Desc.Ar,
						En: menu.Desc.En,
					},
					Image:    menu.Image,
					MobileId: item.MobileId,
					HasError: validationErrorMap[item.Id],
					PriceSummary: responses.PriceSummary{
						Qty:                      item.Qty,
						UnitPrice:                menu.Price,
						TotalPriceBeforeDiscount: SubTotalMenusValueBefore,
						//seperate it after apply offer
						TotalPriceAfterDiscount: SubTotalMenusValueAfter,
					},
					ModifierItems: modifierDocs,
				}
				menuRespDocs = append(menuRespDocs, menuRespDoc)
			}
		}
	}
	return
}

func GetModifierItemPriceSummary(item order.MenuItem, menu extMenuResponses.MenuDetailsResponse, validationErrorMap map[string][]string) (modifierRespDocs []responses.MenuDoc, totalModifierValueBefore float64, totalModifierValueAfter float64) {
	modifierRespDocs = make([]responses.MenuDoc, 0)
	for _, modifier := range item.ModifierIds {
		for _, addon := range menu.Addons {
			if modifier.Id == addon.ID.Hex() {
				SubTotalModifierValueBefore := float64(modifier.Qty) * addon.Price
				SubTotalModifierValueAfter := float64(modifier.Qty) * addon.Price
				totalModifierValueBefore += SubTotalModifierValueBefore
				totalModifierValueAfter += SubTotalModifierValueAfter

				modifierRespDoc := responses.MenuDoc{
					Id: modifier.Id,
					Name: responses.LocalizationText{
						Ar: addon.Name.Ar,
						En: addon.Name.En,
					},
					Image:    addon.Image,
					HasError: validationErrorMap[item.Id],
					PriceSummary: responses.PriceSummary{
						Qty:                      modifier.Qty,
						UnitPrice:                addon.Price,
						TotalPriceBeforeDiscount: SubTotalModifierValueBefore,
						//seperate it after apply offer
						TotalPriceAfterDiscount: SubTotalModifierValueAfter,
					},
				}
				modifierRespDocs = append(modifierRespDocs, modifierRespDoc)
			}
		}
	}
	return
}

func CheckIsLocationReadyForNewOrder(ctx *context.Context, doc extLocResponses.LocationDetails) validators.ErrorResponse {
	if !doc.IsOpen {
		return validators.GetErrorResponse(ctx, localization.Mobile_location_not_open_error, nil, utils.GetAsPointer(http.StatusUnprocessableEntity))
	}
	return validators.ErrorResponse{}
}

func CheckIsMenuItemsValid(ctx *context.Context, menuDocs []extMenuResponses.MenuDetailsResponse, menuItemDto []order.MenuItem) validators.ErrorResponse {
	modifierMap := make(map[string]bool)
	menuItemsMap := make(map[string]bool)

	for _, menu := range menuDocs {
		menuItemsMap[menu.ID.Hex()] = true
		for _, addon := range menu.Addons {
			modifierMap[addon.ID.Hex()] = true
		}
	}

	ValidationErrors := make(map[string][]string)
	for menuIndex, item := range menuItemDto {
		if _, ok := menuItemsMap[item.Id]; !ok {
			ValidationErrors[fmt.Sprintf("menu.%d", menuIndex)] = []string{localization.GetTranslation(ctx, localization.Mobile_item_unavailable, nil, "")}
		}
		for modifierIndex, modifier := range item.ModifierIds {
			if _, ok := modifierMap[modifier.Id]; !ok {
				ValidationErrors[fmt.Sprintf("menu.%d.modifier.%d", menuIndex, modifierIndex)] = []string{localization.GetTranslation(ctx, localization.Mobile_item_unavailable, nil, "")}
			}
		}
	}

	if len(ValidationErrors) > 1 {
		return validators.ErrorResponse{IsError: true, ValidationErrors: ValidationErrors, StatusCode: http.StatusUnprocessableEntity}
	}

	return validators.ErrorResponse{}
}
