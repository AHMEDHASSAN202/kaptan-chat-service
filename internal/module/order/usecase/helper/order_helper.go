package helper

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/asaskevich/EventBus"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	extMenuResponses "samm/internal/module/order/external/menu/responses"
	extLocResponses "samm/internal/module/order/external/retails/responses"
	"samm/internal/module/order/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
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

func GenerateSerialNumber() string {
	currentTime := time.Now().Format("020106")
	randomNumber := rand.Intn(900000) + 10000
	serialNumber := fmt.Sprintf("%s-%06d", currentTime, randomNumber)
	return serialNumber
}
func GetNextAndPreviousStatusByType(actor string, currentStatus string, nextStatus string) (nextStatuses []string, previousStatus []string) {

	var orderStatus map[string]domain.OrderStatusJson
	dir, err := os.Getwd()
	if err != nil {
		return nextStatuses, previousStatus
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "order_status.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nextStatuses, previousStatus
	}

	if errRe := json.Unmarshal(data, &orderStatus); errRe != nil {
		return nextStatuses, previousStatus
	}
	statusRule := orderStatus[currentStatus]
	nextStatusRule := orderStatus[nextStatus]

	switch actor {
	case consts.ActorAdmin:
		return statusRule.AllowAdminToChange, nextStatusRule.PreviousStatus
	case consts.ActorKitchen:
		return statusRule.AllowKitchenToChange, nextStatusRule.PreviousStatus
	case consts.ActorUser:
		return statusRule.AllowUserToChange, nextStatusRule.PreviousStatus
	}
	return nextStatuses, previousStatus

}

func KitchenRejectionReasons(ctx context.Context, status string, id string) ([]domain.KitchenRejectionReason, validators.ErrorResponse) {
	kitchenRejectionReason := make([]domain.KitchenRejectionReason, 0)
	dir, err := os.Getwd()
	if err != nil {
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "kitchen_cancel_reasons.json")
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error(err)
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &kitchenRejectionReason); errRe != nil {
		logger.Logger.Error(err)
		return kitchenRejectionReason, validators.GetErrorResponseFromErr(errRe)
	}

	// Handle Status
	if status != "" {
		From(kitchenRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.KitchenRejectionReason).Status == status || c.(domain.KitchenRejectionReason).Status == "all"
		}).ToSlice(&kitchenRejectionReason)
	}
	if id != "" {
		From(kitchenRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.KitchenRejectionReason).Id == id
		}).ToSlice(&kitchenRejectionReason)
	}

	return kitchenRejectionReason, validators.ErrorResponse{}
}
func UserRejectionReasons(ctx context.Context, status string, id string) ([]domain.UserRejectionReason, validators.ErrorResponse) {
	userRejectionReason := make([]domain.UserRejectionReason, 0)
	dir, err := os.Getwd()
	if err != nil {
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "user_cancel_reasons.json")
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error("Read Json File -> Error -> ", err)
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &userRejectionReason); errRe != nil {
		logger.Logger.Error("ListPermissions -> Error -> ", errRe)
		return userRejectionReason, validators.GetErrorResponseFromErr(errRe)
	}

	// Handle Status
	if status != "" {
		From(userRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.UserRejectionReason).Status == status || c.(domain.UserRejectionReason).Status == "all"
		}).ToSlice(&userRejectionReason)
	}
	if id != "" {
		From(userRejectionReason).Where(func(c interface{}) bool {
			return c.(domain.UserRejectionReason).Id == id
		}).ToSlice(&userRejectionReason)
	}

	return userRejectionReason, validators.ErrorResponse{}
}

func NotifyUser(order *domain.Order, bus EventBus.Bus) {
	//// Send notification
	//var notificationData notification.NotificationDto
	//notificationData.Title.Ar = ""
	//notificationData.Title.En = ""
	//notificationData.Text.En = ""
	//notificationData.Text.Ar = ""
	//notificationData.Type = consts3.TYPE_PRIVATE
	//notificationData.Ids = []string{utils.ConvertObjectIdToStringId(order.User.ID)}
	//notificationData.ModelType = consts2.UserModelType
	//notificationData.CountryId = order.User.Country

	//bus.Publish(consts3.SEND_NOTIFICATION, notificationData)
}
