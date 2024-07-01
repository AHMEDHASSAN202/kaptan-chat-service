package admin

import (
	admin2 "samm/internal/module/admin/builder/role"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func AdminProfileBuilder(model *domain.Admin) *admin.AdminProfileResponse {
	return &admin.AdminProfileResponse{
		ID:         model.ID,
		Name:       model.Name,
		Email:      model.Email,
		Type:       model.Type,
		Role:       admin2.FindRoleBuilder(&model.Role),
		CountryIds: model.CountryIds,
		AccountId:  model.MetaData.AccountId,
	}
}
