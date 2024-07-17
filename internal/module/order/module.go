package order

import (
	"go.uber.org/fx"
	"samm/internal/module/order/delivery"
	"samm/internal/module/order/external"
	"samm/internal/module/order/policies"
	order_repo "samm/internal/module/order/repository/order"
	order_usecase "samm/internal/module/order/usecase/order"
	"samm/internal/module/order/usecase/order_factory"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		order_repo.NewOrderMongoRepository,
		order_factory.NewOrderFactory,
		order_usecase.NewOrderUseCase,
		external.NewExternalService,
	),
	fx.Invoke(
		policies.NewIPolicy, delivery.InitOrderController,
	),
)
