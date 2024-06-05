package database

import (
	"go.uber.org/fx"
	"samm/pkg/database/mongodb"
)

var Module = fx.Options(
	fx.Provide(
		mongodb.NewClient,
	),
)
