package admin

import (
	"go.uber.org/fx"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/delivery"
	admin2 "samm/internal/module/admin/repository/mongodb/admin"
	"samm/internal/module/admin/usecase/admin"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		admin2.NewAdminRepository,
		admin.NewAdminUseCase,
		custom_validators.InitNewCustomValidatorsAdmin,
	),
	fx.Invoke(
		delivery.InitAdminController,
	),
)
