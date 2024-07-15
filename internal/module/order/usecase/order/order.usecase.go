package order

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	"net/http"
	"os"
	"path/filepath"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external"
	"samm/internal/module/order/responses"
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
func (l OrderUseCase) ListOrderForDashboard(ctx context.Context, payload *order.ListOrderDto) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListOrderForDashboard(&ctx, payload)
	if dbErr != nil {
		return nil, validators.GetErrorResponseFromErr(dbErr)
	}
	return responses.SetListResponse(ordersRes, paginationMeta), validators.ErrorResponse{}
}

func (l OrderUseCase) StoreOrder(ctx context.Context, payload *order.CreateOrderDto) (interface{}, validators.ErrorResponse) {
	orderFactory, err := l.orderFactory.Make("ktha")
	if err != nil {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

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
	hasLocErr := helper.CheckIsLocationReadyForNewOrder(&ctx, locationDoc)
	if hasLocErr.IsError {
		l.logger.Error(hasLocErr.ErrorMessageObject.Text)
		return resp, hasLocErr
	}
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

func (l OrderUseCase) UserRejectionReasons(ctx context.Context, status string, id string) ([]domain.UserRejectionReason, validators.ErrorResponse) {
	userRejectionReason := make([]domain.UserRejectionReason, 0)
	dir, err := os.Getwd()
	if err != nil {
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	path := filepath.Join(dir, "internal", "module", "order", "assets", "user_cancel_reasons.json")
	data, err := os.ReadFile(path)
	if err != nil {
		l.logger.Error("Read Json File -> Error -> ", err)
		return userRejectionReason, validators.GetErrorResponseFromErr(err)
	}

	if errRe := json.Unmarshal(data, &userRejectionReason); errRe != nil {
		l.logger.Error("ListPermissions -> Error -> ", errRe)
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

func (l OrderUseCase) UserCancelOrder(ctx context.Context, payload *order.CancelOrderDto) (*domain.Order, validators.ErrorResponse) {
	// Find Order
	orderDomain, err := l.repo.FindOrder(&ctx, payload.OrderId, payload.UserId)
	if err != nil {
		return orderDomain, validators.GetErrorResponseFromErr(err)
	}
	// Check Status
	nextStatuses, previousStatuses := helper.GetNextAndPreviousStatusByType(consts.ActorUser, orderDomain.Status, consts.OrderStatus.Cancelled)
	if !utils.Contains(nextStatuses, consts.OrderStatus.Cancelled) {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}

	rejectionReason, errRe := l.UserRejectionReasons(ctx, "", payload.CancelReasonId)
	if errRe.IsError {
		return nil, errRe
	}
	if len(rejectionReason) == 0 {
		return nil, validators.GetErrorResponseWithErrors(&ctx, localization.ChangeOrderStatusError, nil)
	}
	now := time.Now().UTC()
	updateSet := map[string]interface{}{
		"status":       consts.OrderStatus.Cancelled,
		"cancelled_at": now,
		"updated_at":   now,
		"cancelled": domain.Rejected{
			Id:       payload.CancelReasonId,
			Note:     "",
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

	orderDomain, err = l.repo.UpdateOrderStatus(&ctx, orderDomain, previousStatuses, statusLog, updateSet)

	fmt.Println("err")
	fmt.Println(err)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	// If Status Is Pending Call Payment To Release This Transaction

	// Send Notification

	return orderDomain, validators.ErrorResponse{}
}
