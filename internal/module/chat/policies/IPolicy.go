package policies

import (
	"kaptan/internal/module/chat/domain"
	"kaptan/pkg/gate"
)

type IPolicy struct {
}

func NewIPolicy(g *gate.Gate) *IPolicy {
	g.Register(&domain.Chat{}, &ChatPolicy{})
	g.Register(&domain.Message{}, &MessagePolicy{})
	return &IPolicy{}
}
