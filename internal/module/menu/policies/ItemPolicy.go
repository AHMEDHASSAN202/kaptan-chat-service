package policies

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

type ItemPolicy struct {
}

func (p ItemPolicy) Before(item *domain.Item, ctx context.Context) bool {
	if ctx.Value("causer-type") == "admin" {
		return true
	}
	return false
}

func (p ItemPolicy) Check(item *domain.Item, ctx context.Context) bool {
	if p.Before(item, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(item.AccountId) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p ItemPolicy) Find(item *domain.Item, ctx context.Context) bool {
	return p.Check(item, ctx)
}

func (p ItemPolicy) Update(item *domain.Item, ctx context.Context) bool {
	return p.Check(item, ctx)
}

func (p ItemPolicy) Delete(item *domain.Item, ctx context.Context) bool {
	return p.Check(item, ctx)
}
