package user

import (
	"samm/internal/module/user/domain"
	"samm/pkg/jwt"
	"samm/pkg/logger"
)

type Middlewares struct {
	userRepository domain.UserRepository
	logger         logger.ILogger
	jwtFactory     jwt.JwtServiceFactory
}

func NewUserMiddlewares(userRepository domain.UserRepository, logger logger.ILogger, jwtFactory jwt.JwtServiceFactory) *Middlewares {
	return &Middlewares{
		userRepository: userRepository,
		logger:         logger,
		jwtFactory:     jwtFactory,
	}
}
