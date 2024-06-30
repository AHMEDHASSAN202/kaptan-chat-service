package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func AdminProfileBuilder(model *domain.Admin) *admin.AdminProfileResponse {
	return &admin.AdminProfileResponse{
		ID:          model.ID,
		Name:        model.Name,
		Email:       model.Email,
		Type:        model.Type,
		Role:        model.Role,
		Permissions: model.Permissions,
		CountryIds:  model.CountryIds,
		AccountId:   model.MetaData.AccountId,
	}
}
