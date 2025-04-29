package middlewares

import (
	"go.uber.org/fx"
	commmon "kaptan/pkg/middlewares/common"
	"kaptan/pkg/middlewares/user"
)

var Module = fx.Options(
	fx.Provide(
		user.NewUserMiddlewares,
		commmon.NewCommonMiddlewares,
	),
)
