package user

import (
	"samm/internal/module/user/domain"
	"samm/pkg/database/redis"
	"samm/pkg/jwt"
	"samm/pkg/logger"
)

type Middlewares struct {
	userRepository domain.UserRepository
	redisClient    *redis.RedisClient
	logger         logger.ILogger
	jwtFactory     jwt.JwtServiceFactory
}

func NewUserMiddlewares(userRepository domain.UserRepository, logger logger.ILogger, jwtFactory jwt.JwtServiceFactory, redisClient *redis.RedisClient) *Middlewares {
	return &Middlewares{
		userRepository: userRepository,
		redisClient:    redisClient,
		logger:         logger,
		jwtFactory:     jwtFactory,
	}
}
