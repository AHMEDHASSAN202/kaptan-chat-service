package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func ListAdminBuilder(models *[]domain.Admin) *[]admin.FindAdminResponse {
	data := make([]admin.FindAdminResponse, 0)
	if models == nil {
		return &data
	}
	for _, model := range *models {
		ad := FindAdminBuilder(&model)
		data = append(data, *ad)
	}
	return &data
}
