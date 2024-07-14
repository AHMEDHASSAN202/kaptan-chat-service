package user

import (
	"context"
	"github.com/jinzhu/copier"
	"net/http"
	"samm/internal/module/order/domain"
	"samm/internal/module/order/responses/user"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func FindOrderBuilder(ctx *context.Context, orderModel *domain.Order) (*user.FindOrderResponse, validators.ErrorResponse) {
	orderResponse := user.FindOrderResponse{}
	err := copier.Copy(&orderResponse, orderModel)
	if err != nil {
		return nil, validators.GetErrorResponse(ctx, localization.CanNotBuildOrderResponse, nil, utils.GetAsPointer(http.StatusInternalServerError))
	}
	return &orderResponse, validators.ErrorResponse{}
}
