package user

import (
	"go.uber.org/fx"
	"kaptan/internal/module/user/migrations"
	"kaptan/internal/module/user/repository/driver"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		driver.NewDriverRepository,
	),
	migrations.Module,
)
