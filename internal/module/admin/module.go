package admin

import (
	"go.uber.org/fx"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/delivery"
	admin2 "samm/internal/module/admin/repository/mongodb/admin"
	"samm/internal/module/admin/repository/mongodb/role"
	"samm/internal/module/admin/usecase/admin"
	role2 "samm/internal/module/admin/usecase/role"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		admin2.NewAdminRepository,
		role.NewRoleRepository,
		admin.NewAdminUseCase,
		role2.NewRoleUseCase,
		custom_validators.InitNewCustomValidatorsAdmin,
		custom_validators.InitNewCustomValidatorsRole,
	),
	fx.Invoke(
		delivery.InitAdminController, delivery.InitAdminAuthController, delivery.InitAdminPortalController, delivery.InitRoleController,
	),
)
