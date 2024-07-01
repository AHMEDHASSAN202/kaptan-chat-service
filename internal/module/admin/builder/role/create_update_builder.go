package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/role"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
)

func CreateUpdateRoleBuilder(role *domain.Role, input *dto.CreateRoleDTO) (*domain.Role, error) {
	if role == nil {
		role = &domain.Role{}
		role.ID = primitive.NewObjectID()
		role.AdminDetails = []dto2.AdminDetails{}
	}
	role.Name.En = input.Name.En
	role.Name.Ar = input.Name.Ar
	role.Permissions = utils.ArrayToLower(input.Permissions)
	role.AdminDetails = append(role.AdminDetails, input.AdminDetails)
	return role, nil
}
