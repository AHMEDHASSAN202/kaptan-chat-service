package kitchen

import (
	"samm/internal/module/admin/domain"
	"samm/pkg/database/redis"
	"samm/pkg/jwt"
	"samm/pkg/logger"
)

type ProviderMiddlewares struct {
	adminRepository domain.AdminRepository
	logger          logger.ILogger
	jwtFactory      jwt.JwtServiceFactory
	redisClient     *redis.RedisClient
}

func NewKitchenMiddlewares(adminRepository domain.AdminRepository, logger logger.ILogger, jwtFactory jwt.JwtServiceFactory, redisClient *redis.RedisClient) *ProviderMiddlewares {
	return &ProviderMiddlewares{
		adminRepository: adminRepository,
		logger:          logger,
		jwtFactory:      jwtFactory,
		redisClient:     redisClient,
	}
}
