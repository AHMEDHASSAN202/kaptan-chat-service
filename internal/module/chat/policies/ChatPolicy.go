package policies

import (
	"context"
	"kaptan/internal/module/chat/domain"
)

type ChatPolicy struct {
}

func (p ChatPolicy) Before(admin *domain.Chat, ctx context.Context) bool {
	return false
}

func (p ChatPolicy) Update(order *domain.Chat, ctx context.Context) bool {
	return false
}

func (p ChatPolicy) Delete(order *domain.Chat, ctx context.Context) bool {
	return false
}
