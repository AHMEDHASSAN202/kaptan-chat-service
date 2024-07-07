package policies

import (
	"samm/internal/module/admin/domain"
	"samm/pkg/gate"
)

type IPolicy struct {
}

func NewIPolicy(g *gate.Gate) *IPolicy {
	g.Register(&domain.Admin{}, &AdminPolicy{})
	return &IPolicy{}
}
