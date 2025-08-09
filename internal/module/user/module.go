package user

import (
	"go.uber.org/fx"
	"kaptan/internal/module/user/migrations"
	"kaptan/internal/module/user/repository/driver"
	"kaptan/internal/module/user/usecase/notification"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		driver.NewDriverRepository,
		notification.NewNotificationUseCase,
	),
	migrations.Module,
)
