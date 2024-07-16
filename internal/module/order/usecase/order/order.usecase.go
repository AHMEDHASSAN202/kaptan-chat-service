package order

import (
	"context"
	"errors"
	"net/http"
	user2 "samm/internal/module/order/builder/user"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external"
	"samm/internal/module/order/responses"
	"samm/internal/module/order/responses/user"
	"samm/internal/module/order/usecase/helper"
	"samm/internal/module/order/usecase/order_factory"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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

func (l OrderUseCase) ListOrderForMobile(ctx context.Context, payload *order.ListOrderDtoForMobile) (*responses.ListResponse, validators.ErrorResponse) {
	ordersRes, paginationMeta, dbErr := l.repo.ListOrderForMobile(&ctx, payload)
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
	orderDomain, errRe := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(payload.OrderId))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}

	if orderDomain.User.ID.Hex() != payload.UserId {
		l.logger.Error(" User >> unauthorized access ")
		return validators.GetErrorResponse(ctx, localization.E1006, nil, nil)
	}

	if orderDomain.IsFavourite {
		orderDomain.IsFavourite = false
	} else {
		orderDomain.IsFavourite = true
	}
	errRe = l.repo.UpdateOrder(orderDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
