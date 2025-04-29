package policies

import (
	"context"
	"fmt"
	"kaptan/internal/module/chat/domain"
)

type MessagePolicy struct {
}

func (p MessagePolicy) Check(message *domain.Message, ctx context.Context) bool {
	if ctx.Value("causer-id") == fmt.Sprintf("%s", message.SenderId) && ctx.Value("causer-type") == message.SenderType {
		return true
	}
	return false
}

func (p MessagePolicy) Update(message *domain.Message, ctx context.Context) bool {
	return p.Check(message, ctx)
}

func (p MessagePolicy) Delete(message *domain.Message, ctx context.Context) bool {
	return p.Check(message, ctx)
}
