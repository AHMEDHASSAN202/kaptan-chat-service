package policies

import (
	"samm/internal/module/menu/domain"
	"samm/pkg/gate"
)

type IPolicy struct {
}

func NewIPolicy(g *gate.Gate) *IPolicy {
	g.Register(&domain.MenuGroup{}, &MenuGroupPolicy{})
	return &IPolicy{}
}
