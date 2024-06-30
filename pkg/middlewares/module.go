package middlewares

import (
	"go.uber.org/fx"
	"samm/pkg/middlewares/admin"
	"samm/pkg/middlewares/portal"
)

var Module = fx.Options(
	fx.Provide(
		admin.NewAdminMiddlewares,
		portal.NewPortalMiddlewares,
	),
)
