package database

import (
	"example.com/fxdemo/pkg/database/mongodb"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		mongodb.NewClient,
	),
)
