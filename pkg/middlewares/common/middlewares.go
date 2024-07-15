package commmon

import (
	"samm/pkg/database/redis"
	"samm/pkg/logger"
)

type ProviderMiddlewares struct {
	logger      logger.ILogger
	RedisClient *redis.RedisClient
}

func NewCommonMiddlewares(logger logger.ILogger, redisClient *redis.RedisClient) *ProviderMiddlewares {
	return &ProviderMiddlewares{
		logger:      logger,
		RedisClient: redisClient,
	}
}
