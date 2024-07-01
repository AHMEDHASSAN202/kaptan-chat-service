package admin

import (
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/role"
)

func FindRoleBuilder(model *domain.Role) role.FindRoleResponse {
	return role.FindRoleResponse{
		ID:          model.ID,
		Name:        role.Name{En: model.Name.En, Ar: model.Name.Ar},
		Permissions: model.Permissions,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
