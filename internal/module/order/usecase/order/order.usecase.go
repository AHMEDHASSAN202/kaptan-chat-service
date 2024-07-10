package order

import (
	"context"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external"
	"samm/internal/module/order/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
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

func (l OrderUseCase) StoreOrder(ctx context.Context, payload *order.StoreOrderDto) (err validators.ErrorResponse) {
	orderDomain := domain.Order{}
	orderDomain.Name.Ar = payload.Name.Ar
	orderDomain.Name.En = payload.Name.En
	orderDomain.Email = payload.Email
	password, er := utils.HashPassword(payload.Password)
	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	orderDomain.Password = password
	orderDomain.CreatedAt = time.Now()
	orderDomain.UpdatedAt = time.Now()

	errRe := l.repo.StoreOrder(ctx, &orderDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l OrderUseCase) UpdateOrder(ctx context.Context, id string, payload *order.UpdateOrderDto) (err validators.ErrorResponse) {
	orderDomain, errRe := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	orderDomain.Name.Ar = payload.Name.Ar
	orderDomain.Name.En = payload.Name.En
	orderDomain.Email = payload.Email

	if payload.Password != "" {
		password, er := utils.HashPassword(payload.Password)
		if er != nil {
			return validators.GetErrorResponseFromErr(er)
		}
		orderDomain.Password = password
	}
	orderDomain.UpdatedAt = time.Now()

	errRe = l.repo.UpdateOrder(ctx, orderDomain)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l OrderUseCase) FindOrder(ctx context.Context, Id string) (order domain.Order, err validators.ErrorResponse) {
	domainOrder, errRe := l.repo.FindOrder(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainOrder, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainOrder, validators.ErrorResponse{}
}

func (l OrderUseCase) DeleteOrder(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteOrder(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l OrderUseCase) ListOrder(ctx context.Context, payload *order.ListOrderDto) (orders []domain.Order, paginationResult utils.PaginationResult, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListOrder(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}

func (l OrderUseCase) CalculateOrderCost(ctx context.Context, payload *order.CalculateOrderCostDto) (resp responses.CalculateOrderCostResp, err validators.ErrorResponse) {
	//find location details
	locationDoc, errResponse := l.extService.RetailsIService.GetLocationDetails(ctx, payload.LocationId)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.Mobile_location_not_available_error, nil, nil)
	}
	//check is the location available for the order
	hasLocErr := checkIsLocationReadyForNewOrder(&ctx, locationDoc)
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
	resp, errResponse = l.calculateOrderCostBuilder(ctx, locationDoc, menuDetails, payload)
	if errResponse.IsError {
		l.logger.Error(errResponse.ErrorMessageObject.Text)
		return resp, validators.GetErrorResponse(&ctx, localization.E1005, nil, nil)
	}
	return resp, validators.ErrorResponse{}
}
