package database

import (
	"go.uber.org/fx"
	"samm/pkg/database/mongodb"
	"samm/pkg/database/redis"
)

var Module = fx.Options(
	fx.Provide(
		mongodb.NewClient,
		redis.NewRedisClient,
	),
)
