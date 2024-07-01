package admin

import (
	admin2 "samm/internal/module/admin/builder/role"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func FindAdminBuilder(model *domain.Admin) *admin.FindAdminResponse {
	return &admin.FindAdminResponse{
		ID:         model.ID,
		Name:       model.Name,
		Email:      model.Email,
		Type:       model.Type,
		Role:       admin2.FindRoleBuilder(&model.Role),
		CountryIds: model.CountryIds,
		MetaData:   admin.MetaData{AccountId: model.MetaData.AccountId},
		Status:     model.Status,
		CreatedAt:  model.CreatedAt,
		UpdateAt:   model.UpdatedAt,
	}
}
