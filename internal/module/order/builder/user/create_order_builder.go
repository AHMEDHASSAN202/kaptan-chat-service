package user

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	responses2 "samm/internal/module/order/external/menu/responses"
	"samm/internal/module/order/external/retails/responses"
	"samm/internal/module/order/usecase/helper"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

// CreateOrderBuilder that creates a new order based on the provided data.
// It also includes a helper function buildItemsAndSummary that populates the order with items and calculates the order summary.
func CreateOrderBuilder(ctx context.Context, dto *order.CreateOrderDto, location responses.LocationDetails, items []responses2.MenuDetailsResponse, collectionMethod *responses.CollectionMethod, accountDoc responses.AccountDetails) (*domain.Order, validators.ErrorResponse) {
	orderModel := domain.Order{}
	orderModel.ID = primitive.NewObjectID()
	orderModel.CreatedAt = time.Now().UTC()
	orderModel.UpdatedAt = time.Now().UTC()
	orderModel.SerialNum = helper.GenerateSerialNumber()
	orderModel.User = domain.User{}
	orderModel.Items = []domain.Item{}
	buildItemsAndSummary(&orderModel, dto, items)
	if err := copier.Copy(&orderModel.Location, location); err != nil {
		logger.Logger.Error(err)
	}
	if err := copier.Copy(&orderModel.Location.Brand, location.BrandDetails); err != nil {
		logger.Logger.Error(err)
	}
	if err := copier.Copy(&orderModel.User, ctx.Value("causer-details")); err != nil {
		logger.Logger.Error(err)
	}
	if collectionMethod != nil {
		var collectionMethodObj domain.CollectionMethod
		if err := copier.Copy(&collectionMethodObj, collectionMethod); err != nil {
			logger.Logger.Error(err)
		} else {
			orderModel.User.CollectionMethod = &collectionMethodObj
		}
	}
	if err := copier.Copy(&orderModel.Location.Account, accountDoc); err != nil {
		logger.Logger.Error(err)
	}
	orderModel.PreparationTime = location.PreparationTime
	orderModel.Status = consts.OrderStatus.Initiated
	orderModel.IsFavourite = dto.IsFavourite
	orderModel.StatusLogs = []domain.StatusLog{}
	orderModel.Notes = dto.Notes
	return &orderModel, validators.ErrorResponse{}
}

// buildItemsAndSummary populates the order with items and calculates the order summary
func buildItemsAndSummary(order *domain.Order, dto *order.CreateOrderDto, items []responses2.MenuDetailsResponse) {
	var totalMenusValueBefore float64
	var totalMenusValueAfter float64
	// Map of items and their modifiers
	itemsMap := map[string]responses2.MenuDetailsResponse{}
	modifiersMaps := map[string]map[string]responses2.MobileGetItemAddon{}

	// Populate the items map and modifiers map
	for _, item := range items {
		itemsMap[utils.ConvertObjectIdToStringId(item.ID)] = item
		if item.Addons != nil {
			modifiersMaps[utils.ConvertObjectIdToStringId(item.ID)] = map[string]responses2.MobileGetItemAddon{}
			for _, addon := range item.Addons {
				modifiersMaps[utils.ConvertObjectIdToStringId(item.ID)][utils.ConvertObjectIdToStringId(addon.ID)] = addon
			}
		}
	}

	// Build the items for the order
	for _, item := range dto.MenuItems {
		menuItem := itemsMap[item.Id]
		menuItemOrder := domain.Item{
			ID:       menuItem.ID,
			ItemId:   menuItem.ItemId,
			MobileId: item.MobileId,
			Name: domain.LocalizationText{
				En: menuItem.Name.En,
				Ar: menuItem.Name.Ar,
			},
			Desc: domain.LocalizationText{
				En: menuItem.Desc.En,
				Ar: menuItem.Desc.Ar,
			},
			Min:      0,
			Max:      0,
			Calories: menuItem.Calories,
			Price:    menuItem.Price,
			Image:    menuItem.Image,
			Qty:      int(item.Qty),
			Addons:   []domain.Item{},
		}

		// Calculate total values
		modifierDocs, totalModifierValueBefore, totalModifierValueAfter := helper.GetModifierItemPriceSummary(item, menuItem, make(map[string][]string))

		// Add the modifiers (addons) to the menuItemOrder
		for _, addon := range item.ModifierIds {
			modifier := modifiersMaps[item.Id][addon.Id]
			addonItem := domain.Item{
				ID: modifier.ID,
				Name: domain.LocalizationText{
					En: modifier.Name.En,
					Ar: modifier.Name.Ar,
				},
				Min:             modifier.Min,
				Max:             modifier.Max,
				Calories:        modifier.Calories,
				Price:           modifier.Price,
				Image:           modifier.Image,
				Qty:             int(addon.Qty),
				ModifierGroupId: getModifierGroupId(menuItem.ModifierGroups, modifier.ID),
			}

			for _, doc := range modifierDocs {
				if doc.Id == addon.Id {
					addonItem.PriceSummary = domain.ItemPriceSummary{
						Qty:                      int(doc.PriceSummary.Qty),
						UnitPrice:                doc.PriceSummary.UnitPrice,
						TotalPriceBeforeDiscount: doc.PriceSummary.TotalPriceBeforeDiscount,
						TotalPriceAfterDiscount:  doc.PriceSummary.TotalPriceAfterDiscount,
					}
				}
			}

			menuItemOrder.Addons = append(menuItemOrder.Addons, addonItem)
		}

		SubTotalMenusValueBefore := float64(item.Qty) * (menuItem.Price + totalModifierValueBefore)
		SubTotalMenusValueAfter := float64(item.Qty) * (menuItem.Price + totalModifierValueAfter)

		menuItemOrder.PriceSummary = domain.ItemPriceSummary{
			Qty:                      int(item.Qty),
			UnitPrice:                menuItem.Price,
			TotalPriceBeforeDiscount: SubTotalMenusValueBefore,
			TotalPriceAfterDiscount:  SubTotalMenusValueAfter,
		}

		totalMenusValueBefore += SubTotalMenusValueBefore
		totalMenusValueAfter += SubTotalMenusValueAfter

		order.Items = append(order.Items, menuItemOrder)
	}

	order.PriceSummary = domain.OrderPriceSummary{
		Fees:                     0,
		TotalPriceBeforeDiscount: totalMenusValueBefore,
		TotalPriceAfterDiscount:  totalMenusValueAfter,
	}
}

func getModifierGroupId(modifierGroups []responses2.ModifierGroup, addonId primitive.ObjectID) *primitive.ObjectID {
	if modifierGroups == nil {
		return nil
	}
	for _, group := range modifierGroups {
		if utils.Contains(group.ProductIds, addonId) {
			return &group.ID
		}
	}
	return nil
}
