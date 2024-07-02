package middlewares

import (
	"go.uber.org/fx"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/middlewares/portal"
	"samm/pkg/middlewares/user"
)

var Module = fx.Options(
	fx.Provide(
		admin.NewAdminMiddlewares,
		portal.NewPortalMiddlewares,
		user.NewUserMiddlewares,
		commmon.NewCommonMiddlewares,
	),
)
