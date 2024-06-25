package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func ListAdminBuilder(models *[]domain.Admin) *[]admin.ListAdminResponse {
	data := make([]admin.ListAdminResponse, 0)
	if models == nil {
		return &data
	}
	for _, model := range *models {
		data = append(data, admin.ListAdminResponse{
			ID:         model.ID,
			Name:       model.Name,
			Email:      model.Email,
			Type:       model.Type,
			Role:       model.Role,
			CountryIds: model.CountryIds,
			MetaData:   admin.MetaData{AccountId: model.MetaData.AccountId},
			Status:     model.Status,
			CreatedAt:  model.CreatedAt,
			UpdateAt:   model.UpdatedAt,
		})
	}
	return &data
}
