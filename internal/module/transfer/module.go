package transfer

import (
	"go.uber.org/fx"
	"kaptan/internal/module/transfer/repository"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		repository.NewTransferRepository,
	),
)
