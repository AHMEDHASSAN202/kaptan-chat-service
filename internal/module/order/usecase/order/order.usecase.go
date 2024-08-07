package order

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/db"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	user2 "samm/internal/module/order/builder/user"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/dto/order/kitchen"
	"samm/internal/module/order/external"
	"samm/internal/module/order/responses"
	"samm/internal/module/order/responses/user"
	"samm/internal/module/order/usecase/helper"
	"samm/internal/module/order/usecase/order_factory"
	"samm/pkg/gate"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strings"
	"sync"
	"time"
)

type OrderUseCase struct {
	repo         domain.OrderRepository
	extService   external.ExtService
	logger       logger.ILogger
	orderFactory *order_factory.OrderFactory
	realTimeDb   *db.Client
	gate         *gate.Gate
}

func NewOrderUseCase(repo domain.OrderRepository, extService external.ExtService, gate *gate.Gate, realTimeDb *db.Client, logger logger.ILogger, orderFactory *order_factory.OrderFactory) domain.OrderUseCase {
	return &OrderUseCase{
		repo:         repo,
		extService:   extService,
		logger:       logger,
		orderFactory: orderFactory,
		realTimeDb:   realTimeDb,
		gate:         gate,
	}
}
func (l OrderUseCase) ListOrderForDashboard(ctx context.Context, payload *order.ListOrderDtoForDashboard) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListOrderForDashboard(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}

func (l OrderUseCase) ListInprogressOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListInprogressOrdersForMobile(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}

func (l OrderUseCase) ListCompletedOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListCompletedOrdersForMobile(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}

func (l OrderUseCase) ListLastOrdersForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListLastOrdersForMobile(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}

func (l *OrderUseCase) FindOrderForDashboard(ctx *context.Context, id string) (*domain.Order, validators.ErrorResponse) {
	order, err := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if order == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return order, validators.ErrorResponse{}
}

func (l *OrderUseCase) FindOrderForMobile(ctx *context.Context, payload *order.FindOrderMobileDto) (orderResponse *user.FindOrderResponse, err validators.ErrorResponse) {
	order, dbErr := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if dbErr != nil {
		err = validators.GetErrorResponseFromErr(dbErr)
		return
	}
	if order == nil {
		err = validators.GetErrorResponseFromErr(errors.New(localization.E1002))
		return
	}

	if order.User.ID.Hex() != payload.UserId {
		l.logger.Error(" User >> unauthorized access ")
		err = validators.GetErrorResponse(ctx, localization.E1006, nil, nil)
		return
	}

	//builder order response
	orderResponse, err = user2.FindOrderBuilder(ctx, order)

	return
}

func (l OrderUseCase) StoreOrder(ctx context.Context, payload *order.CreateOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//create order
	orderResponse, errCreate := orderFactory.Create(ctx, payload)
	if errCreate.IsError {
		return nil, errCreate
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	//find location details
	locationDoc, errResponse := l.extService.RetailsIService.GetLocationDetails(ctx, payload.LocationId)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.Mobile_location_not_available_error, nil, nil)
	}
	//check is the location available for the order
	//hasLocErr := helper.CheckIsLocationReadyForNewOrder(&ctx, locationDoc)
	//if hasLocErr.IsError {
	//	l.logger.Error(hasLocErr.ErrorMessageObject.Text)
	//	return resp, hasLocErr
	//}
	//find menus details
	menuDetails, errResponse := l.extService.MenuIService.GetMenuItemsDetails(ctx, payload.MenuItems, payload.LocationId)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.E1002Item, nil, nil)
	}
	//check is the menus are available
	resp, errResponse = helper.CalculateOrderCostBuilder(ctx, locationDoc, menuDetails, payload)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.E1005, nil, nil)
	}
	return resp, validators.ErrorResponse{}
}

func (l OrderUseCase) ToggleOrderFavourite(ctx *context.Context, payload order.ToggleOrderFavDto) (err validators.ErrorResponse) {
	orderDomain, dbErr := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}

	if orderDomain.User.ID.Hex() != payload.UserId {
		l.logger.Error(" User >> unauthorized access ")
		return validators.GetErrorResponse(ctx, localization.E1006, nil, nil)
	}

	if orderDomain.IsFavourite {
		orderDomain.IsFavourite = false
		dbErr = l.repo.UpdateOrder(*ctx, orderDomain)
		if dbErr != nil {
			return validators.GetErrorResponseFromErr(dbErr)
		}
	} else {
		orderDomain.IsFavourite = true
		transactionErr := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
			dbErr = l.repo.UpdateUserAllOrdersFavorite(sc, orderDomain.User.ID.Hex())
			if dbErr != nil {
				return dbErr
			}
			dbErr = l.repo.UpdateOrder(sc, orderDomain)
			if dbErr != nil {
				return dbErr
			}
			return session.CommitTransaction(sc)
		})
		if transactionErr != nil {
			return validators.GetErrorResponseFromErr(transactionErr)
		}
	}
	return
}

func (l OrderUseCase) UserRejectionReasons(ctx context.Context, status string, id string) ([]domain.UserRejectionReason, validators.ErrorResponse) {
	return helper.UserRejectionReasons(ctx, status, id)
}
func (l OrderUseCase) ReportMissedItem(ctx context.Context, payload *order.ReportMissingItemDto) (interface{}, validators.ErrorResponse) {
	orderDoc, err := l.repo.FindOrder(&ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if err != nil {
		l.logger.Error("ReportMissedItem", err.Error())
		return nil, validators.GetErrorResponseFromErr(err)
	}

	//authorize this order
	if !l.gate.Authorize(orderDoc, "ReportMissingItem", ctx) {
		l.logger.Error("ReportMissingItem -> UnAuthorized update -> ", orderDoc.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	//check is reported before
	if orderDoc.MetaData.HasMissingItems {
		l.logger.Error("ReportMissingItem -> reported before-> ", orderDoc.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.Reported_Item_Already_Added, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	missingItemMap := make(map[string]order.MissedItems)
	missingAddonMap := make(map[string]map[string]order.MissedItems)
	for _, item := range payload.MissingItems {
		missingItemMap[item.Id] = item
		for _, addon := range item.MissingAddons {
			if missingAddonMap[item.Id] == nil {
				missingAddonMap[item.Id] = make(map[string]order.MissedItems)
			}
			missingAddonMap[item.Id][addon.Id] = addon
		}
	}

	for i, item := range orderDoc.Items {
		if val, ok := missingItemMap[item.MobileId]; ok {
			if orderDoc.Items[i].MissedItemReport == nil {
				orderDoc.Items[i].MissedItemReport = &domain.MissedItem{}
			}
			orderDoc.Items[i].MissedItemReport.Id = val.Id
			orderDoc.Items[i].MissedItemReport.Qty = val.Qty
			orderDoc.MetaData.HasMissingItems = true
			for addonIndex, addon := range item.Addons {
				if val, ok := missingAddonMap[item.MobileId][addon.ID.Hex()]; ok {
					if orderDoc.Items[i].Addons[addonIndex].MissedItemReport == nil {
						orderDoc.Items[i].Addons[addonIndex].MissedItemReport = &domain.MissedItem{}
					}
					orderDoc.Items[i].Addons[addonIndex].MissedItemReport.Id = val.Id
					orderDoc.Items[i].Addons[addonIndex].MissedItemReport.Qty = val.Qty
				}
			}
		}
	}
	err = l.repo.UpdateOrder(ctx, orderDoc)
	if err != nil {
		l.logger.Error("ReportMissedItem", err)
		return nil, validators.GetErrorResponseFromErr(err)
	}

	return nil, validators.ErrorResponse{}
}

func (l OrderUseCase) UserCancelOrder(ctx context.Context, payload *order.CancelOrderDto) (*user.FindOrderResponse, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToCancel(ctx, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	return orderResponse, validators.ErrorResponse{}

}
func (l OrderUseCase) DashboardCancelOrder(ctx context.Context, payload *order.DashboardCancelOrderDto) (*domain.Order, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToCancelDashboard(ctx, consts.ActorAdmin, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	return orderResponse, validators.ErrorResponse{}

}
func (l OrderUseCase) DashboardPickedOrder(ctx context.Context, payload *order.DashboardPickedUpOrderDto) (*domain.Order, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToPickedUpDashboard(ctx, consts.ActorAdmin, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	return orderResponse, validators.ErrorResponse{}

}

func (l OrderUseCase) UserArrivedOrder(ctx context.Context, payload *order.ArrivedOrderDto) (*user.FindOrderResponse, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToArrived(ctx, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) SetOrderPaid(ctx context.Context, payload *order.OrderPaidDto) validators.ErrorResponse {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	errPaid := orderFactory.ToPaid(ctx, payload)
	if errPaid.IsError {
		return errPaid
	}
	return validators.ErrorResponse{}

}

func (l OrderUseCase) KitchenAcceptOrder(ctx context.Context, payload *kitchen.AcceptOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToAcceptKitchen(ctx, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) KitchenRejectedOrder(ctx context.Context, payload *kitchen.RejectedOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errRejected := orderFactory.ToRejectedKitchen(ctx, payload)
	if errRejected.IsError {
		return nil, errRejected
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) KitchenRejectionReasons(ctx context.Context, status string, id string) ([]domain.KitchenRejectionReason, validators.ErrorResponse) {
	return helper.KitchenRejectionReasons(ctx, status, id)
}

func (l OrderUseCase) KitchenReadyForPickupOrder(ctx context.Context, payload *kitchen.ReadyForPickupOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errRejected := orderFactory.ToReadyForPickupKitchen(ctx, payload)
	if errRejected.IsError {
		return nil, errRejected
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) KitchenPickedUpOrder(ctx context.Context, payload *kitchen.PickedUpOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errRejected := orderFactory.ToPickedUpKitchen(ctx, payload)
	if errRejected.IsError {
		return nil, errRejected
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) KitchenNoShowOrder(ctx context.Context, payload *kitchen.NoShowOrderDto) (interface{}, validators.ErrorResponse) {
	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errRejected := orderFactory.ToNoShowKitchen(ctx, payload)
	if errRejected.IsError {
		return nil, errRejected
	}

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) UpdateRealTimeDb(ctx context.Context, order *domain.Order) validators.ErrorResponse {
	ttlArr := []string{}
	err := l.realTimeDb.NewRef("orders/raw").Child(order.ID.Hex()).Set(ctx, utils.ObjectToStringified(order))
	if err != nil {
		l.logger.Error("raw: sync to firebase order => ", err)
		return validators.GetErrorResponseFromErr(err)
	}
	ttlArr = append(ttlArr, "^orderId", order.ID.Hex())

	err = l.realTimeDb.NewRef("orders/users").Child(order.User.ID.Hex()).Child(order.ID.Hex()).Set(ctx, order.Status)
	if err != nil {
		l.logger.Error("users: sync to firebase order => ", err)
		return validators.GetErrorResponseFromErr(err)
	}
	ttlArr = append(ttlArr, "^userId", order.User.ID.Hex())

	for _, kitchenId := range order.MetaData.TargetKitchenIds {
		err = l.realTimeDb.NewRef("orders/kitchens").Child(kitchenId.Hex()).Child(order.ID.Hex()).Set(ctx, order.Status)
		if err != nil {
			l.logger.Error("kitchens: sync to firebase order => ", err)
			return validators.GetErrorResponseFromErr(err)
		}
		ttlArr = append(ttlArr, "^kitchenId", kitchenId.Hex())
	}

	//setup ttl field with its expired at
	err = l.realTimeDb.NewRef("orders/ttl").Child(strings.Join(ttlArr, " ")).Set(ctx, time.Now().UTC().Add(6*time.Hour).Format(utils.DefaultDateTimeFormat))
	if err != nil {
		l.logger.Error("users: sync to firebase order => ", err)
		return validators.GetErrorResponseFromErr(err)
	}

	return validators.ErrorResponse{}
}

// cron jobs

func (l OrderUseCase) CronJobTimedOutOrders(ctx context.Context) validators.ErrorResponse {
	// last 5 hours interval
	t := time.Now().UTC().Add(-5 * time.Hour)
	tRFC, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
		bson.M{"status": consts.OrderStatus.Pending},
		bson.M{"created_at": bson.M{"$lte": tRFC}}},
	}}

	orders, _, dbErr := l.repo.GetAllOrdersForCronJobs(&ctx, matching)
	if dbErr != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	if orders == nil {
		return validators.ErrorResponse{}
	}

	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}
	errChan := make(chan validators.ErrorResponse, len(*orders))
	var wg sync.WaitGroup
	for _, v := range *orders {
		wg.Add(1)
		go func(order domain.Order) validators.ErrorResponse {
			defer wg.Done()
			toTimedOutErr := orderFactory.ToTimedOut(ctx, order.ID.Hex())
			if toTimedOutErr.IsError {
				l.logger.Error("Error processing  TimedOut order:", toTimedOutErr)
				errChan <- toTimedOutErr
			}
			return validators.ErrorResponse{}
		}(v)
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return <-errChan
	}

	return validators.ErrorResponse{}
}

func (l OrderUseCase) CronJobPickedOrders(ctx context.Context) validators.ErrorResponse {
	// last 5 hours interval
	t := time.Now().UTC().Add(-5 * time.Hour)
	tRFC, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
		bson.M{"status": consts.OrderStatus.Accepted},
		bson.M{"created_at": bson.M{"$lte": tRFC}}},
	}}

	orders, _, dbErr := l.repo.GetAllOrdersForCronJobs(&ctx, matching)
	if dbErr != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	if orders == nil {
		return validators.ErrorResponse{}
	}

	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	errChan := make(chan validators.ErrorResponse, len(*orders))
	var wg sync.WaitGroup
	for _, v := range *orders {
		wg.Add(1)
		go func(orderData domain.Order) validators.ErrorResponse {
			defer wg.Done()
			payload := &order.DashboardPickedUpOrderDto{
				OrderId: orderData.ID.Hex(),
				AdminHeaders: dto.AdminHeaders{
					CauserType: "cron",
				},
			}
			_, pickedErr := orderFactory.ToPickedUpDashboard(ctx, consts.ActorCron, payload)
			if pickedErr.IsError {
				l.logger.Error("Error processing PickedUp order:", pickedErr)
				errChan <- pickedErr
			}
			return validators.ErrorResponse{}
		}(v)
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return <-errChan
	}

	return validators.ErrorResponse{}
}

func (l OrderUseCase) CronJobCancelOrders(ctx context.Context) validators.ErrorResponse {
	t := time.Now().UTC().Add(-5 * time.Hour)
	tRFC, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
		bson.M{"status": consts.OrderStatus.Initiated},
		bson.M{"created_at": bson.M{"$lte": tRFC}}},
	}}

	orders, _, dbErr := l.repo.GetAllOrdersForCronJobs(&ctx, matching)
	if dbErr != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	if orders == nil {
		return validators.ErrorResponse{}
	}

	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	errChan := make(chan validators.ErrorResponse, len(*orders))
	var wg sync.WaitGroup
	for _, v := range *orders {
		wg.Add(1)
		go func(orderData domain.Order) validators.ErrorResponse {
			defer wg.Done()
			payload := &order.DashboardCancelOrderDto{
				OrderId: orderData.ID.Hex(),
				AdminHeaders: dto.AdminHeaders{
					CauserType: "cron",
				},
			}
			_, cancelErr := orderFactory.ToCancelDashboard(ctx, consts.ActorCron, payload)
			if cancelErr.IsError {
				l.logger.Error("Error processing Cancelled order:", cancelErr)
				errChan <- cancelErr
			}

			return validators.ErrorResponse{}
		}(v)
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return <-errChan
	}

	return validators.ErrorResponse{}
}

func (l OrderUseCase) KitchenListRunningOrders(ctx context.Context, payload *kitchen.ListRunningOrderDto) (interface{}, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListRunningOrdersForKitchen(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}
