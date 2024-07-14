package order_factory

import (
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
	"samm/internal/module/order/builder"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/dto/order"
	"samm/internal/module/order/external"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

//OrderFactory.Make("ktha")->Create()->sendNotification();
//OrderFactory.Make("ktha")->Find(id);
//OrderFactory.Make("ktha")->List(dto);
//OrderFactory.Make("ktha")->Find(id)->ToPending(dto);
//OrderFactory.Make("ktha")->Find(id)->ToAccept(dto);
//OrderFactory.Make("ktha")->Find(id)->ToArrived(dto);
//OrderFactory.Make("ktha"->Find(id))->ToCancel(dto);

type Deps struct {
	validator  *validator.Validate
	extService *external.ExtService
	logger     logger.ILogger
	orderRepo  domain.OrderRepository
}

type NgoOrder struct {
	Order domain.Order
	Deps
}

func (o *NgoOrder) Make() IOrder {
	return o
}

func (o *NgoOrder) Create(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
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
		return nil, validators.GetErrorResponse(&ctx, localization.Mobile_location_not_available_error, nil, utils.GetAsPointer(http.StatusUnprocessableEntity))
	}

	//check is the location available for the order
	hasLocErr := CheckIsLocationReadyForNewOrder(&ctx, locationDoc)
	if hasLocErr.IsError {
		o.logger.Error(hasLocErr.ErrorMessageObject.Text)
		return nil, hasLocErr
	}

	//find menus details
	menuDetails, errResponse := o.extService.MenuIService.GetMenuItemsDetails(ctx, input.MenuItems, input.LocationId)
	if errResponse.IsError {
		o.logger.Error(errResponse.ErrorMessageObject.Text)
		return nil, validators.GetErrorResponse(&ctx, localization.E1002Item, nil, nil)
	}

	//check menu items are available for the order
	hasMenuErr := CheckIsMenuItemsValid(&ctx, menuDetails, input.MenuItems)
	if hasMenuErr.IsError {
		o.logger.Error(hasMenuErr.ErrorMessageObject.Text)
		return nil, hasMenuErr
	}

	//order builder
	orderModel, errOrderModel := builder.CreateOrderBuilder(ctx, input, locationDoc, menuDetails)
	if errOrderModel.IsError {
		o.logger.Error(hasMenuErr.ErrorMessageObject.Text)
		return nil, hasMenuErr
	}

	errStoreOrder := o.orderRepo.StoreOrder(&ctx, orderModel)
	if errStoreOrder != nil {
		o.logger.Error(errStoreOrder)
		return nil, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}

	return nil, validators.ErrorResponse{}
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

func (o *NgoOrder) ToAccept(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) ToArrived(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}

func (o *NgoOrder) ToCancel(ctx context.Context, dto interface{}) (*domain.Order, validators.ErrorResponse) {
	return nil, validators.ErrorResponse{}
}
