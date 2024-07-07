package policies

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

type MenuGroupPolicy struct {
}

func (p MenuGroupPolicy) Before(menuGroup domain.MenuGroup, ctx context.Context) bool {
	if ctx.Value("causer-type") == "admin" {
		return true
	}
	return false
}

func (p MenuGroupPolicy) Check(menuGroup domain.MenuGroup, ctx context.Context) bool {
	if p.Before(menuGroup, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(menuGroup.AccountId) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p MenuGroupPolicy) Find(menuGroup domain.MenuGroup, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}

func (p MenuGroupPolicy) Update(menuGroup domain.MenuGroup, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}

func (p MenuGroupPolicy) Delete(menuGroup domain.MenuGroup, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}
