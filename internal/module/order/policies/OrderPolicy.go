package policies

import (
	"context"
	"samm/internal/module/order/domain"
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
	return true
}

func (p OrderPolicy) KitchenToRejected(order *domain.Order, ctx context.Context) bool {
	return true
}
