package database

import (
	"go.uber.org/fx"
	"kaptan/pkg/database/mysql"
)

var Module = fx.Options(
	fx.Provide(
		mysql.NewClient,
	),
)
