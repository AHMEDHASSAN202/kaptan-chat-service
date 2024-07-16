package policies

import (
	"samm/internal/module/order/domain"
	"samm/pkg/gate"
)

type IPolicy struct {
}

func NewIPolicy(g *gate.Gate) *IPolicy {
	g.Register(&domain.Order{}, &OrderPolicy{})
	return &IPolicy{}
}
