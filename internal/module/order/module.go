package order

import (
	"go.uber.org/fx"
	"samm/internal/module/order/delivery"
	order_repo "samm/internal/module/order/repository/order"
	order_usecase "samm/internal/module/order/usecase/order"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		order_repo.NewOrderMongoRepository,
		order_usecase.NewOrderUseCase,
	),
	fx.Invoke(
		delivery.InitOrderController,
	),
)
