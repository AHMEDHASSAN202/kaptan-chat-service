package middlewares

import (
	"go.uber.org/fx"
	"samm/pkg/middlewares/admin"
	"samm/pkg/middlewares/portal"
	"samm/pkg/middlewares/user"
)

var Module = fx.Options(
	fx.Provide(
		admin.NewAdminMiddlewares,
		portal.NewPortalMiddlewares,
		user.NewUserMiddlewares,
	),
)
