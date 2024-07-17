package policies

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

type ModifierGroupPolicy struct {
}

func (p ModifierGroupPolicy) Before(modifier *domain.ModifierGroup, ctx context.Context) bool {
	if ctx.Value("causer-type") == "admin" {
		return true
	}
	return false
}

func (p ModifierGroupPolicy) Check(modifier *domain.ModifierGroup, ctx context.Context) bool {
	if p.Before(modifier, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(modifier.AccountId) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p ModifierGroupPolicy) Find(modifier *domain.ModifierGroup, ctx context.Context) bool {
	return p.Check(modifier, ctx)
}

func (p ModifierGroupPolicy) Update(modifier *domain.ModifierGroup, ctx context.Context) bool {
	return p.Check(modifier, ctx)
}

func (p ModifierGroupPolicy) Delete(modifier *domain.ModifierGroup, ctx context.Context) bool {
	return p.Check(modifier, ctx)
}
