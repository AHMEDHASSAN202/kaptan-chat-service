package policies

import (
	"samm/internal/module/menu/domain"
	"samm/pkg/gate"
)

type IPolicy struct {
}

func NewIPolicy(g *gate.Gate) *IPolicy {
	g.Register(&domain.MenuGroup{}, &MenuGroupPolicy{})
	g.Register(&domain.Item{}, &ItemPolicy{})
	g.Register(&domain.ModifierGroup{}, &ModifierGroupPolicy{})
	return &IPolicy{}
}
