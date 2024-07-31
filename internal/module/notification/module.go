package notification

import (
	"go.uber.org/fx"
	domain2 "samm/internal/module/kitchen/domain"
	"samm/internal/module/notification/delivery"
	"samm/internal/module/notification/external"
	"samm/internal/module/notification/gateways/onesignal"
	notification_repo "samm/internal/module/notification/repository/notification"
	notification_usecase "samm/internal/module/notification/usecase/notification"
	"samm/internal/module/user/domain"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		external.NewExternalService,
		notification_repo.NewNotificationMongoRepository,
		notification_usecase.NewNotificationUseCase,
		onesignal.NewOnesignalService,
	),
	fx.Invoke(
		delivery.InitNotificationController,
		func(a *external.ExtService, userUseCase domain.UserUseCase, kitchenUserCase domain2.KitchenUseCase) {
			a.SetUserUseCase(userUseCase)
			a.SetKitchenUseCase(kitchenUserCase)
		},
	),
)
