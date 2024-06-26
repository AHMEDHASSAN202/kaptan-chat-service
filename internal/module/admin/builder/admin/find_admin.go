package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func FindAdminBuilder(model *domain.Admin) *admin.FindAdminResponse {
	return &admin.FindAdminResponse{
		ID:          model.ID,
		Name:        model.Name,
		Email:       model.Email,
		Type:        model.Type,
		Role:        model.Role,
		Permissions: model.Permissions,
		CountryIds:  model.CountryIds,
		MetaData:    admin.MetaData{AccountId: model.MetaData.AccountId},
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		UpdateAt:    model.UpdatedAt,
	}
}
