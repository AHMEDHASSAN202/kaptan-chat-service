package payment

import (
	"go.uber.org/fx"
	"samm/internal/module/payment/delivery"
	"samm/internal/module/payment/gateways/myfatoorah"
	"samm/internal/module/payment/repository/mongodb"
	"samm/internal/module/payment/usecase/card"
	"samm/internal/module/payment/usecase/payment"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		// App Config
		mongodb.NewCardMongoRepository,
		mongodb.NewPaymentMongoRepository,
		payment.NewPaymentUseCase,
		myfatoorah.NewMyFatoorahService,
		card.NewCardUseCase,
	),
	fx.Invoke(
		delivery.InitCardController,
		delivery.InitPaymentController,
	),
)
