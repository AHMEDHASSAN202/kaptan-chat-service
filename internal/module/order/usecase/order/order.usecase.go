package order

import (
	"context"
	"errors"
	"github.com/kamva/mgm/v3"
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
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type OrderUseCase struct {
	repo         domain.OrderRepository
	extService   external.ExtService
	logger       logger.ILogger
	orderFactory *order_factory.OrderFactory
}

func NewOrderUseCase(repo domain.OrderRepository, extService external.ExtService, logger logger.ILogger, orderFactory *order_factory.OrderFactory) domain.OrderUseCase {
	return &OrderUseCase{
		repo:         repo,
		extService:   extService,
		logger:       logger,
		orderFactory: orderFactory,
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

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

	return orderResponse, validators.ErrorResponse{}

}
func (l OrderUseCase) DashboardCancelOrder(ctx context.Context, payload *order.DashboardCancelOrderDto) (*domain.Order, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToCancelDashboard(ctx, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	//send notifications
	go orderFactory.SendNotifications()

	return orderResponse, validators.ErrorResponse{}

}
func (l OrderUseCase) DashboardPickedOrder(ctx context.Context, payload *order.DashboardPickedUpOrderDto) (*domain.Order, validators.ErrorResponse) {

	//create new instance from ktha factory
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//accept order
	orderResponse, errAccept := orderFactory.ToPickedUpDashboard(ctx, payload)
	if errAccept.IsError {
		return nil, errAccept
	}

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

	return orderResponse, validators.ErrorResponse{}
}

func (l OrderUseCase) SetOrderPaid(ctx context.Context, payload *order.OrderPaidDto) validators.ErrorResponse {
	// Find Order
	orderDomain, err := l.repo.FindOrder(&ctx, payload.OrderId)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorUser, orderDomain.Status, consts.OrderStatus.Pending)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Cancelled) {
		return validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":  consts.OrderStatus.Pending,
		"paid_at": now,
		"cancelled": domain.Payment{
			Id:          utils.ConvertStringIdToObjectId(payload.TransactionId),
			PaymentType: payload.PaymentType,
			CardType:    payload.CardType,
			CardNumber:  payload.CardNumber,
		},
	}

	statusLog := domain.StatusLog{
		CauserId:   utils.ConvertObjectIdToStringId(orderDomain.User.ID),
		CauserType: "user",
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.Pending
	statusLog.Status.Old = orderDomain.Status

	orderDomain, err = l.repo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)

	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	// Send Notification

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

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

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

	//send notifications
	go orderFactory.SendNotifications()

	return orderResponse, validators.ErrorResponse{}
}
