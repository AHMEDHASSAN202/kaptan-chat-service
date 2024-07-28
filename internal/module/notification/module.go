package notification

import (
	"go.uber.org/fx"
	"samm/internal/module/notification/delivery"
	notification_repo "samm/internal/module/notification/repository/notification"
	notification_usecase "samm/internal/module/notification/usecase/notification"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		notification_repo.NewNotificationMongoRepository,
		notification_usecase.NewNotificationUseCase,
	),
	fx.Invoke(
		delivery.InitNotificationController,
	),
)
