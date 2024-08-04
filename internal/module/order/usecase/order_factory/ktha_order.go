package order_factory

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"net/http"
	kitchen2 "samm/internal/module/order/builder/kitchen"
	user2 "samm/internal/module/order/builder/user"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/dto/order/kitchen"
	"samm/internal/module/order/external"
	"samm/internal/module/order/external/retails/responses"
	kitchen3 "samm/internal/module/order/responses/kitchen"
	"samm/internal/module/order/responses/user"
	"samm/internal/module/order/subscribers"
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
	eventBus    EventBus.Bus
}

type KthaOrder struct {
	Order *domain.Order
	Deps
}

func (o *KthaOrder) Make() IOrder {
	return o
}

func (o *KthaOrder) Create(ctx context.Context, dto interface{}) (*user.FindOrderResponse, validators.ErrorResponse) {
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
	accountDoc, errResponse := o.extService.RetailsIService.GetAccountDetails(ctx, utils.ConvertObjectIdToStringId(locationDoc.AccountId))
	if errResponse.IsError {
		o.logger.Error(errResponse.ErrorMessageObject.Text)
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.Mobile_location_not_available_error, nil)
	}

	kitchenIds, errResponse := o.extService.KitchenIService.GetKitchensForSpecificLocation(ctx, input.LocationId, utils.ConvertObjectIdToStringId(locationDoc.AccountId))
	if errResponse.IsError {
		o.logger.Error(errResponse.ErrorMessageObject.Text)
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
	var collectionMethod *responses.CollectionMethod
	if input.CollectionMethodId != "" {
		collectionMethodResp, hasCollectionMethodErr := o.extService.RetailsIService.FindCollectionMethod(ctx, input.CollectionMethodId, ctx.Value("causer-id").(string))
		if hasCollectionMethodErr.IsError {
			o.logger.Error(hasCollectionMethodErr)
			return nil, validators.GetErrorResponseWithErrors(&ctx, localization.OrderCollectionMethodError, nil)
		}
		collectionMethod = &collectionMethodResp
	}

	//order builder
	orderModel, errOrderModel := user2.CreateOrderBuilder(ctx, input, locationDoc, menuDetails, collectionMethod, accountDoc, kitchenIds)
	if errOrderModel.IsError {
		o.logger.Error(errOrderModel.ErrorMessageObject.Text)
		return nil, errOrderModel
	}

	//check if user has running orders
	hasRunningOrders, errHasRunningOrders := o.orderRepo.UserHasOrders(ctx, orderModel.User.ID, []string{consts.OrderStatus.Initiated}, 4)
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

	//push an event
	go o.PushEventToSubscribers(ctx)

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) Find(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *KthaOrder) List(ctx context.Context, dto interface{}) ([]domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *KthaOrder) SendNotifications() validators.ErrorResponse {

	return validators.ErrorResponse{}
}

func (o *KthaOrder) ToPending(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *KthaOrder) ToAcceptKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
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

	//authorize order amount
	_, paymentError := o.extService.PaymentIService.AuthorizePayment(ctx, utils.ConvertObjectIdToStringId(orderDomain.Payment.Id), true)
	if paymentError.IsError {
		o.logger.Error("PAYMENT_ERROR => ", paymentError)
		//return nil, validators.GetErrorResponseWithErrors(&ctx, localization.PaymentError, nil)
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

	//push an event
	go o.PushEventToSubscribers(ctx)

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) ToRejectedKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
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

	//authorize order amount
	_, paymentError := o.extService.PaymentIService.AuthorizePayment(ctx, utils.ConvertObjectIdToStringId(orderDomain.Payment.Id), true)
	if paymentError.IsError {
		o.logger.Error("PAYMENT_ERROR => ", paymentError)
		//return nil, validators.GetErrorResponseWithErrors(&ctx, localization.PaymentError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":      consts.OrderStatus.Rejected,
		"rejected_at": now,
		"updated_at":  now,
		"rejected": domain.Rejected{
			Id:   input.RejectedReasonId,
			Note: input.Note,
			Name: &domain.Name{
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

	//push an event
	go o.PushEventToSubscribers(ctx)

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) ToReadyForPickupKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*kitchen.ReadyForPickupOrderDto)
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
	if !o.gate.Authorize(orderDomain, "KitchenToReadyForPickup", ctx) {
		o.logger.Error("ToAcceptKitchen -> UnAuthorized Accept -> ", orderDomain.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorKitchen, orderDomain.Status, consts.OrderStatus.ReadyForPickup)
	if !utils.Contains(nextStatuses, consts.OrderStatus.ReadyForPickup) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":              consts.OrderStatus.ReadyForPickup,
		"ready_for_pickup_at": now,
		"updated_at":          now,
	}
	statusLog := domain.StatusLog{
		CauserId:   input.CauserId,
		CauserType: input.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.ReadyForPickup
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) ToPickedUpKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*kitchen.PickedUpOrderDto)
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
	if !o.gate.Authorize(orderDomain, "KitchenToPickedUp", ctx) {
		o.logger.Error("ToAcceptKitchen -> UnAuthorized Accept -> ", orderDomain.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorKitchen, orderDomain.Status, consts.OrderStatus.PickedUp)
	if !utils.Contains(nextStatuses, consts.OrderStatus.PickedUp) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":      consts.OrderStatus.PickedUp,
		"pickedup_at": now,
		"updated_at":  now,
	}
	statusLog := domain.StatusLog{
		CauserId:   input.CauserId,
		CauserType: input.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.PickedUp
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) ToNoShowKitchen(ctx context.Context, dto interface{}) (*kitchen3.FindOrderResponse, validators.ErrorResponse) {
	//validate
	input := dto.(*kitchen.NoShowOrderDto)
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
	if !o.gate.Authorize(orderDomain, "KitchenToNoShow", ctx) {
		o.logger.Error("ToAcceptKitchen -> UnAuthorized Accept -> ", orderDomain.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorKitchen, orderDomain.Status, consts.OrderStatus.NoShow)
	if !utils.Contains(nextStatuses, consts.OrderStatus.NoShow) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":     consts.OrderStatus.NoShow,
		"no_show_at": now,
		"updated_at": now,
	}
	statusLog := domain.StatusLog{
		CauserId:   input.CauserId,
		CauserType: input.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.NoShow
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}

	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	//build response
	orderResponse, errBuilder := kitchen2.FindOrderBuilder(&ctx, orderDomain)
	if errBuilder.IsError {
		o.logger.Error(errBuilder)
		return nil, errBuilder
	}

	return orderResponse, validators.ErrorResponse{}
}

func (o *KthaOrder) ToArrived(ctx context.Context, payload *order.ArrivedOrderDto) (*user.FindOrderResponse, validators.ErrorResponse) {
	// Find Order
	orderDomain, err := o.orderRepo.FindOrderByUser(&ctx, payload.OrderId, payload.UserId)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	if !utils.Contains([]string{consts.OrderStatus.Accepted, consts.OrderStatus.ReadyForPickup, consts.OrderStatus.NoShow}, orderDomain.Status) || orderDomain.ArrivedAt != nil {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"arrived_at": now,
		"updated_at": now,
	}

	collectionMethod, hasCollectionMethodErr := o.extService.RetailsIService.FindCollectionMethod(ctx, payload.CollectionMethodId, payload.UserId)
	if hasCollectionMethodErr.IsError {
		o.logger.Error(hasCollectionMethodErr)
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.OrderCollectionMethodError, nil)
	}
	var userCollectionMethod domain.CollectionMethod
	copier.Copy(&userCollectionMethod, collectionMethod)
	updateSet["user.collection_method"] = userCollectionMethod

	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, []string{}, nil, updateSet)

	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	orderResponse, _ := user2.FindOrderBuilder(&ctx, orderDomain)

	return orderResponse, validators.ErrorResponse{}
}
func (o *KthaOrder) ToPaid(ctx context.Context, payload *order.OrderPaidDto) validators.ErrorResponse {

	// Find Order
	orderDomain, err := o.orderRepo.FindOrder(&ctx, payload.OrderId)
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

	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)

	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	o.Order = orderDomain

	go o.PushEventToSubscribers(ctx)

	return validators.ErrorResponse{}

}

func (o *KthaOrder) ToCancel(ctx context.Context, payload *order.CancelOrderDto) (*user.FindOrderResponse, validators.ErrorResponse) {

	// Find Order
	orderDomain, err := o.orderRepo.FindOrderByUser(&ctx, payload.OrderId, payload.UserId)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorUser, orderDomain.Status, consts.OrderStatus.Cancelled)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Cancelled) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	rejectionReasons, errRe := helper.UserRejectionReasons(ctx, "", payload.CancelReasonId)
	if errRe.IsError {
		return nil, errRe
	}
	if len(rejectionReasons) == 0 {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}
	rejectionReason := rejectionReasons[0]
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":       consts.OrderStatus.Cancelled,
		"cancelled_at": now,
		"updated_at":   now,
		"cancelled": domain.Rejected{
			Id:   payload.CancelReasonId,
			Note: payload.Note,
			Name: &domain.Name{
				Ar: rejectionReason.Name.Ar,
				En: rejectionReason.Name.En,
			},
			UserType: payload.CauserType,
		},
	}

	statusLog := domain.StatusLog{
		CauserId:   payload.UserId,
		CauserType: payload.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.Cancelled
	statusLog.Status.Old = orderDomain.Status

	// If Status Is Pending Call Payment To Release This Transaction
	if orderDomain.Status == consts.OrderStatus.Pending {
		o.extService.PaymentIService.AuthorizePayment(ctx, utils.ConvertObjectIdToStringId(orderDomain.Payment.Id), false)
	}

	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)

	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	//set order object to factory
	o.Order = orderDomain
	//push an event
	go o.PushEventToSubscribers(ctx)

	orderResponse, _ := user2.FindOrderBuilder(&ctx, orderDomain)

	return orderResponse, validators.ErrorResponse{}
}
func (o *KthaOrder) ToCancelDashboard(ctx context.Context, payload *order.DashboardCancelOrderDto) (*domain.Order, validators.ErrorResponse) {

	// Find Order
	orderDomain, err := o.orderRepo.FindOrder(&ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorAdmin, orderDomain.Status, consts.OrderStatus.Cancelled)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Cancelled) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":       consts.OrderStatus.Cancelled,
		"cancelled_at": now,
		"updated_at":   now,
		"cancelled": domain.Rejected{
			Id:       "",
			Note:     payload.Note,
			Name:     nil,
			UserType: payload.CauserType,
		},
	}

	statusLog := domain.StatusLog{
		CauserId:   payload.CauserId,
		CauserType: payload.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.Cancelled
	statusLog.Status.Old = orderDomain.Status

	// If Status Is Pending Call Payment To Release This Transaction
	if orderDomain.Status == consts.OrderStatus.Pending {
		o.extService.PaymentIService.AuthorizePayment(ctx, utils.ConvertObjectIdToStringId(orderDomain.Payment.Id), false)
	}

	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	return orderDomain, validators.ErrorResponse{}
}
func (o *KthaOrder) ToPickedUpDashboard(ctx context.Context, payload *order.DashboardPickedUpOrderDto) (*domain.Order, validators.ErrorResponse) {

	//find order
	orderDomain, err := o.orderRepo.FindOrder(&ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorAdmin, orderDomain.Status, consts.OrderStatus.PickedUp)
	if !utils.Contains(nextStatuses, consts.OrderStatus.PickedUp) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	//update data
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":      consts.OrderStatus.PickedUp,
		"pickedup_at": now,
		"updated_at":  now,
	}
	statusLog := domain.StatusLog{
		CauserId:   payload.CauserId,
		CauserType: payload.CauserType,
		CreatedAt:  &now,
	}
	statusLog.Status.New = consts.OrderStatus.PickedUp
	statusLog.Status.Old = orderDomain.Status

	//update domain
	orderDomain, err = o.orderRepo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, &statusLog, updateSet)
	if err != nil {
		o.logger.Error(err)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}
	//set order object to factory
	o.Order = orderDomain

	//push an event
	go o.PushEventToSubscribers(ctx)

	return orderDomain, validators.ErrorResponse{}
}
func (o *KthaOrder) PushEventToSubscribers(ctx context.Context) validators.ErrorResponse {
	o.Deps.eventBus.Publish(subscribers.SubscriberTopics.OrderChange, o.Order)
	return validators.ErrorResponse{}
}
