package helper

import (
	"context"
	"samm/internal/module/order/dto/order"
	extMenuResponses "samm/internal/module/order/external/menu/responses"
	extLocResponses "samm/internal/module/order/external/retails/responses"
	"samm/internal/module/order/responses"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

func CalculateOrderCostBuilder(ctx context.Context, loc extLocResponses.LocationDetails, menus []extMenuResponses.MenuDetailsResponse, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	resp.Location = responses.LocationDoc{
		Id: utils.ConvertObjectIdToStringId(loc.Id),
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
	validationErrorMap := CheckIsMenuItemsValid(ctx, menus, payload.MenuItems, false).ValidationErrors
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
