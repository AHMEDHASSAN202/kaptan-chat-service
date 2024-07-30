package common

import (
	"go.uber.org/fx"
	"samm/internal/module/common/delivery"
	"samm/internal/module/common/repository"
	"samm/internal/module/common/usecase/common"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		repository.NewCommonMongoRepository,
		common.NewCommonUseCase,
	),
	fx.Invoke(delivery.InitCommonController),
)
