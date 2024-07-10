package policies

import (
	"context"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/domain"
	"samm/pkg/utils"
)

type AdminPolicy struct {
}

func (p AdminPolicy) Before(admin *domain.Admin, ctx context.Context) bool {
	if ctx.Value("causer-type") == consts.ADMIN_TYPE {
		return true
	}
	return false
}

func (p AdminPolicy) Check(admin *domain.Admin, ctx context.Context) bool {
	if p.Before(admin, ctx) {
		return true
	}
	if utils.ConvertObjectIdToStringId(admin.ID) == ctx.Value("causer-id") {
		return true
	}
	if utils.ConvertObjectIdToStringId(admin.Account.Id) == ctx.Value("causer-account-id") {
		return true
	}
	return false
}

func (p AdminPolicy) Find(admin *domain.Admin, ctx context.Context) bool {
	return p.Check(admin, ctx)
}

func (p AdminPolicy) Update(admin *domain.Admin, ctx context.Context) bool {
	return p.Check(admin, ctx)
}

func (p AdminPolicy) Delete(admin *domain.Admin, ctx context.Context) bool {
	return p.Check(admin, ctx)
}
