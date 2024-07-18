package policies

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

type MenuGroupItemPolicy struct {
}

func (p MenuGroupItemPolicy) Before(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	if ctx.Value("causer-type") == "admin" {
		return true
	}
	return false
}

func (p MenuGroupItemPolicy) Check(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	if p.Before(menuGroup, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(menuGroup.AccountId) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p MenuGroupItemPolicy) Find(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}

func (p MenuGroupItemPolicy) Update(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}

func (p MenuGroupItemPolicy) Delete(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	return p.Check(menuGroup, ctx)
}

func (p MenuGroupItemPolicy) CheckKitchen(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	if p.Before(menuGroup, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(menuGroup.AccountId) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p MenuGroupItemPolicy) UpdateKitchen(menuGroup *domain.MenuGroupItem, ctx context.Context) bool {
	return p.CheckKitchen(menuGroup, ctx)
}
