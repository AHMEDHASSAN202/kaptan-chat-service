package policies

import (
	"context"
	domain2 "samm/internal/module/admin/domain"
	"samm/internal/module/order/consts"
	"samm/internal/module/order/domain"
	"samm/pkg/utils"
)

type OrderPolicy struct {
}

func (p OrderPolicy) Before(admin *domain.Order, ctx context.Context) bool {
	return true
}

func (p OrderPolicy) Check(order *domain.Order, ctx context.Context) bool {
	if p.Before(order, ctx) {
		return true
	}
	return true
}

func (p OrderPolicy) Find(order *domain.Order, ctx context.Context) bool {
	return p.Check(order, ctx)
}

func (p OrderPolicy) KitchenToAccept(order *domain.Order, ctx context.Context) bool {
	kitchenProfile := ctx.Value("causer-details").(*domain2.Admin)
	if !utils.Contains(kitchenProfile.Kitchen.AllowedStatus, consts.OrderStatus.Accepted) {
		return false
	}
	return utils.Contains(kitchenProfile.Kitchen.AccountIds, order.Location.AccountId) || utils.Contains(kitchenProfile.Kitchen.LocationIds, order.Location.ID)
}

func (p OrderPolicy) KitchenToRejected(order *domain.Order, ctx context.Context) bool {
	kitchenProfile := ctx.Value("causer-details").(*domain2.Admin)
	if !utils.Contains(kitchenProfile.Kitchen.AllowedStatus, consts.OrderStatus.Rejected) {
		return false
	}
	return utils.Contains(kitchenProfile.Kitchen.AccountIds, order.Location.AccountId) || utils.Contains(kitchenProfile.Kitchen.LocationIds, order.Location.ID)
}

func (p OrderPolicy) KitchenToPickedUp(order *domain.Order, ctx context.Context) bool {
	kitchenProfile := ctx.Value("causer-details").(*domain2.Admin)
	if !(utils.Contains(kitchenProfile.Kitchen.AllowedStatus, consts.OrderStatus.PickedUp) || utils.Contains(kitchenProfile.Kitchen.AllowedStatus, "delivered")) {
		return false
	}
	return utils.Contains(kitchenProfile.Kitchen.AccountIds, order.Location.AccountId) || utils.Contains(kitchenProfile.Kitchen.LocationIds, order.Location.ID)
}

func (p OrderPolicy) KitchenToNoShow(order *domain.Order, ctx context.Context) bool {
	kitchenProfile := ctx.Value("causer-details").(*domain2.Admin)
	return utils.Contains(kitchenProfile.Kitchen.AccountIds, order.Location.AccountId) || utils.Contains(kitchenProfile.Kitchen.LocationIds, order.Location.ID)
}

func (p OrderPolicy) KitchenToReadyForPickup(order *domain.Order, ctx context.Context) bool {
	kitchenProfile := ctx.Value("causer-details").(*domain2.Admin)
	return utils.Contains(kitchenProfile.Kitchen.AccountIds, order.Location.AccountId) || utils.Contains(kitchenProfile.Kitchen.LocationIds, order.Location.ID)
}
