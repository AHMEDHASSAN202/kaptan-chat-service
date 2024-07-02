package admin

import (
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/role"
	"samm/pkg/utils"
)

func FindRoleBuilder(model *domain.Role) role.FindRoleResponse {
	return role.FindRoleResponse{
		ID:          model.ID,
		Name:        role.Name{En: model.Name.En, Ar: model.Name.Ar},
		Type:        model.Type,
		CanDelete:   !utils.Contains(consts.PreventDeleteRolesIds, utils.ConvertObjectIdToStringId(model.ID)),
		Permissions: model.Permissions,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
