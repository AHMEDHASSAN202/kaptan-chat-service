package transfer

import (
	"go.uber.org/fx"
	"kaptan/internal/module/transfer/delivery"
	"kaptan/internal/module/transfer/migrations"
	"kaptan/internal/module/transfer/repository"
	"kaptan/internal/module/transfer/usecase/transfer"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		repository.NewTransferRepository,
		transfer.NewTransferUseCase,
	),
	fx.Invoke(
		delivery.InitTransferController,
	),
	migrations.Module,
)
