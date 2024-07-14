package order

import (
	"context"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external"
	"samm/internal/module/order/responses"
	"samm/internal/module/order/usecase/order_factory"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type OrderUseCase struct {
	repo       domain.OrderRepository
	extService external.ExtService
	logger     logger.ILogger
}

const tag = " OrderUseCase "

func NewOrderUseCase(repo domain.OrderRepository, extService external.ExtService, logger logger.ILogger) domain.OrderUseCase {
	return &OrderUseCase{
		repo:       repo,
		extService: extService,
		logger:     logger,
	}
}

func (l OrderUseCase) StoreOrder(ctx context.Context, payload *order.CreateOrderDto) (interface{}, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (l OrderUseCase) CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	//find location details
	locationDoc, errResponse := l.extService.RetailsIService.GetLocationDetails(ctx, payload.LocationId)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.Mobile_location_not_available_error, nil, nil)
	}
	//check is the location available for the order
	hasLocErr := order_factory.CheckIsLocationReadyForNewOrder(&ctx, locationDoc)
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
	resp, errResponse = order_factory.CalculateOrderCostBuilder(ctx, locationDoc, menuDetails, payload)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.E1005, nil, nil)
	}
	return resp, validators.ErrorResponse{}
}
