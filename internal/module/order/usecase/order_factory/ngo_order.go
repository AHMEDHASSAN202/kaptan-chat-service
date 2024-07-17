package order_factory

import (
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
	kitchen2 "samm/internal/module/order/builder/kitchen"
	user2 "samm/internal/module/order/builder/user"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/dto/order/kitchen"
	"samm/internal/module/order/external"
	kitchen3 "samm/internal/module/order/responses/kitchen"
	"samm/internal/module/order/responses/user"
	"samm/internal/module/order/usecase/helper"
	"samm/pkg/database/redis"
	"samm/pkg/gate"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strings"
	"time"
)

//OrderFactory.Make("ktha")->Create()->sendNotification();
//OrderFactory.Make("ktha")->Find(id);
//OrderFactory.Make("ktha")->List(dto);
//OrderFactory.Make("ktha")->Find(id)->ToPending(dto);
//OrderFactory.Make("ktha")->Find(id)->ToAccept(dto);
//OrderFactory.Make("ktha")->Find(id)->ToArrived(dto);
//OrderFactory.Make("ktha"->Find(id))->ToCancel(dto);

type Deps struct {
	validator   *validator.Validate
	extService  external.ExtService
	logger      logger.ILogger
	orderRepo   domain.OrderRepository
	redisClient *redis.RedisClient
	gate        *gate.Gate
}

type NgoOrder struct {
	Order *domain.Order
	Deps
}

func (o *NgoOrder) Make() IOrder {
	return o
}

func (o *NgoOrder) Create(ctx context.Context, dto interface{}) (*user.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*order.CreateOrderDto)
	validateErr := input.Validate(ctx, o.Deps.validator)
	if validateErr.IsError {
		return nil, validateErr
	}

	//find location details
	locationDoc, errResponse := o.extService.RetailsIService.GetLocationDetails(ctx, input.LocationId)
	if errResponse.IsError {
		o.logger.Error(errResponse.ErrorMessageObject.Text)
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.Mobile_location_not_available_error, nil)
	}

	//check is the location available for the order
	hasLocErr := helper.CheckIsLocationReadyForNewOrder(&ctx, locationDoc)
	if hasLocErr.IsError {
		o.logger.Error(hasLocErr.ErrorMessageObject.Text)
		return nil, hasLocErr
	}

	//find menus details
	menuDetails, errResponse := o.extService.MenuIService.GetMenuItemsDetails(ctx, input.MenuItems, input.LocationId)
	if errResponse.IsError {
		o.logger.Error(errResponse.ErrorMessageObject.Text)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//check menu items are available for the order
	hasMenuErr := helper.CheckIsMenuItemsValid(&ctx, menuDetails, input.MenuItems, true)
	if hasMenuErr.IsError {
		o.logger.Error(hasMenuErr)
		return nil, hasMenuErr
	}

	//get user collection method
	collectionMethod, hasCollectionMethodErr := o.extService.RetailsIService.FindCollectionMethod(ctx, input.CollectionMethodId, ctx.Value("causer-id").(string))
	if hasCollectionMethodErr.IsError {
		o.logger.Error(hasCollectionMethodErr)
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.OrderCollectionMethodError, nil)
	}

	//order builder
	orderModel, errOrderModel := user2.CreateOrderBuilder(ctx, input, locationDoc, menuDetails, collectionMethod)
	if errOrderModel.IsError {
		o.logger.Error(errOrderModel.ErrorMessageObject.Text)
		return nil, errOrderModel
	}

	//check if user has running orders
	hasRunningOrders, errHasRunningOrders := o.orderRepo.UserHasOrders(ctx, orderModel.User.ID, []string{consts.OrderStatus.Initiated})
	if errHasRunningOrders != nil {
		o.logger.Error(errHasRunningOrders)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}
	if hasRunningOrders {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.UserHasRunningOrders, nil)
	}

	//new lock to prevent user create to order in the same time
	lockKey := strings.Replace(consts.CREATE_ORDER_LOCK_PREFIX, ":userId", utils.ConvertObjectIdToStringId(orderModel.User.ID), 1)
	if errLock := redis.Lock(ctx, o.redisClient, o.logger, lockKey, time.Second*10); errLock.IsError {
		return nil, errLock
	}

	//save order
	orderModel, errStoreOrder := o.orderRepo.StoreOrder(ctx, orderModel)
	if errStoreOrder != nil {
		o.logger.Error(errStoreOrder)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	//unlock the lock
	redis.UnLock(ctx, o.redisClient, o.logger, lockKey)

	//builder order response
	orderResponse, err := user2.FindOrderBuilder(&ctx, orderModel)
	if err.IsError {
		return nil, err
	}

	//set order object to factory
	o.Order = orderModel

	return orderResponse, validators.ErrorResponse{}
}

func (o *NgoOrder) Find(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) List(ctx context.Context, dto interface{}) ([]domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) SendNotifications() validators.ErrorResponse {
	return validators.ErrorResponse{}
}

func (o *NgoOrder) ToPending(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) ToAcceptKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*kitchen.AcceptOrderDto)
	validateErr := input.Validate(ctx, o.Deps.validator)
	if validateErr.IsError {
		return nil, validateErr
	}

	//find order
	orderDomain, err := o.orderRepo.FindOrder(&ctx, utils.ConvertStringIdToObjectId(input.OrderId))
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	//authorize this order
	if !o.gate.Authorize(orderDomain, "KitchenToAccept", ctx) {
		o.logger.Error("ToAcceptKitchen -> UnAuthorized Accept -> ", orderDomain.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorKitchen, orderDomain.Status, consts.OrderStatus.Accepted)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Accepted) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":           consts.OrderStatus.Accepted,
		"preparation_time": input.PreparationTime,
		"accepted_at":      now,
		"updated_at":       now,
	}
	statusLog := domain.StatusLog{
		CauserId:   input.CauserId,
		CauserType: input.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.Accepted
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//set order object to factory
	o.Order = orderDomain

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *NgoOrder) ToRejectedKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*kitchen.RejectedOrderDto)
	validateErr := input.Validate(ctx, o.Deps.validator)
	if validateErr.IsError {
		return nil, validateErr
	}

	//find order
	orderDomain, err := o.orderRepo.FindOrder(&ctx, utils.ConvertStringIdToObjectId(input.OrderId))
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	//authorize this order
	if !o.gate.Authorize(orderDomain, "KitchenToRejected", ctx) {
		o.logger.Error("ToAcceptKitchen -> UnAuthorized Accept -> ", orderDomain.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorKitchen, orderDomain.Status, consts.OrderStatus.Rejected)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Rejected) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	rejectionReasons, errRe := helper.KitchenRejectionReasons(ctx, "", input.RejectedReasonId)
	if errRe.IsError {
		return nil, errRe
	}
	if len(rejectionReasons) == 0 {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}
	rejectionReason := rejectionReasons[0]

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":      consts.OrderStatus.Rejected,
		"rejected_at": now,
		"updated_at":  now,
		"rejected": domain.Rejected{
			Id:   input.RejectedReasonId,
			Note: input.Note,
			Name: domain.Name{
				Ar: rejectionReason.Name.Ar,
				En: rejectionReason.Name.En,
			},
			UserType: input.CauserType,
		},
	}
	statusLog := domain.StatusLog{
		CauserId:   input.CauserId,
		CauserType: input.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.Rejected
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//set order object to factory
	o.Order = orderDomain

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *NgoOrder) ToArrived(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) ToCancel(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}
