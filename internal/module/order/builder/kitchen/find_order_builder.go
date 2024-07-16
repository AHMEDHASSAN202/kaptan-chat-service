package kitchen

import (
	"context"
	"github.com/jinzhu/copier"
	"net/http"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/responses/kitchen"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func FindOrderBuilder(ctx *context.Context, orderModel *domain.Order) (*kitchen.FindOrderResponse, validators.ErrorResponse) {
	orderResponse := kitchen.FindOrderResponse{}
	err := copier.CopyWithOption(&orderResponse, orderModel, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, validators.GetErrorResponse(ctx, localization.CanNotBuildOrderResponse, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}
	return &orderResponse, validators.ErrorResponse{}
}
