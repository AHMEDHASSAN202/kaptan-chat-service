package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/role"
)

func ListRoleBuilder(models *[]domain.Role) *[]role.FindRoleResponse {
	data := make([]role.FindRoleResponse, 0)
	if models == nil {
		return &data
	}
	for _, model := range *models {
		data = append(data, FindRoleBuilder(&model))
	}
	return &data
}
