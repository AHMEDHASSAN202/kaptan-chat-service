package order

import (
	"context"
	"fmt"
	"net/http"
	"samm/internal/module/order/dto/order"
	extMenuResponses "samm/internal/module/order/external/menu/responses"
	extLocResponses "samm/internal/module/order/external/retails/responses"
	"samm/internal/module/order/responses"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func (l OrderUseCase) calculateOrderCostBuilder(ctx context.Context, loc extLocResponses.LocationDetails, menus []extMenuResponses.MenuDetailsResponse, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	resp.Location = responses.LocationDoc{
		Id: loc.Id,
		Name: responses.LocalizationText{
			Ar: loc.Name.Ar,
			En: loc.Name.En,
		},
	}
	var totalMenusValueBefore, totalMenusValueAfter float64
	resp.MenuItems, totalMenusValueBefore, totalMenusValueAfter = calculateTotalCostForMenus(menus, payload)
	resp.TotalPriceSummary = responses.TotalPriceSummary{
		Fees:                     0,
		TotalPriceBeforeDiscount: totalMenusValueBefore,
		TotalPriceAfterDiscount:  totalMenusValueAfter,
	}

	return resp, err
}

func calculateTotalCostForMenus(menus []extMenuResponses.MenuDetailsResponse, payload *order.CalculateOrderCostDto) (menuRespDocs []responses.MenuDoc, totalMenusValueBefore float64, totalMenusValueAfter float64) {
	menuRespDocs = make([]responses.MenuDoc, 0)
	for _, item := range payload.MenuItems {
		for _, menu := range menus {
			if menu.ID.Hex() == item.Id {
				modifierDocs, totalModifierValueBefore, totalModifierValueAfter := getModifierItemPriceSummary(item, menu)

				//calculate total values
				SubTotalMenusValueBefore := float64(item.Qty)*menu.Price + totalModifierValueBefore
				SubTotalMenusValueAfter := float64(item.Qty)*menu.Price + totalModifierValueAfter

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
					Image: menu.Image,
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

func getModifierItemPriceSummary(item order.MenuItem, menu extMenuResponses.MenuDetailsResponse) (modifierRespDocs []responses.MenuDoc, totalModifierValueBefore float64, totalModifierValueAfter float64) {
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
					Image: addon.Image,
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
func checkIsLocationReadyForNewOrder(ctx *context.Context, doc extLocResponses.LocationDetails) validators.ErrorResponse {
	if !doc.IsOpen {
		return validators.GetErrorResponse(ctx, localization.Mobile_location_not_open_error, nil, nil)
	}
	return validators.ErrorResponse{}
}

func checkIsMenuItemsValid(ctx *context.Context, menuDocs []extMenuResponses.MenuDetailsResponse, menuItemDto []order.MenuItem) validators.ErrorResponse {

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
			ValidationErrors[fmt.Sprintf("menu.%d", menuIndex)] = []string{"not exists"}
		}
		for modifierIndex, modifier := range item.ModifierIds {
			if _, ok := modifierMap[modifier.Id]; !ok {
				ValidationErrors[fmt.Sprintf("menu.%d.modifier.%d", menuIndex, modifierIndex)] = []string{"not exists"}
			}
		}
	}
	fmt.Println(ValidationErrors)
	if len(ValidationErrors) > 0 {
		return validators.ErrorResponse{
			ValidationErrors: ValidationErrors,
			IsError:          true,
			ErrorMessageObject: &validators.Message{
				Text: "validationError",
				Code: localization.E1002,
			},
			StatusCode: http.StatusUnprocessableEntity,
		}
	}
	return validators.ErrorResponse{}
}
