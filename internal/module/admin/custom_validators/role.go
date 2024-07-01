package custom_validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/domain"
	"samm/pkg/utils"
)

type RoleCustomValidator struct {
	adminUseCase domain.AdminUseCase
}

func InitNewCustomValidatorsRole(adminUseCase domain.AdminUseCase) RoleCustomValidator {
	return RoleCustomValidator{
		adminUseCase: adminUseCase,
	}
}

func (i *RoleCustomValidator) ValidateRoleHasAdmins() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		isExists, err := i.adminUseCase.CheckRoleExists(context.Background(), utils.ConvertStringIdToObjectId(val))
		if err.IsError {
			return false
		}
		return !isExists
	}
}

func (i *RoleCustomValidator) ValidateStaticRoles() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(string)
		return !utils.Contains(consts.PreventDeleteRolesIds, val)
	}
}
