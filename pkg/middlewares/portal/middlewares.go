package portal

import (
	"samm/internal/module/admin/domain"
	"samm/pkg/jwt"
	"samm/pkg/logger"
)

type ProviderMiddlewares struct {
	adminRepository domain.AdminRepository
	logger          logger.ILogger
	jwtFactory      jwt.JwtServiceFactory
}

func NewPortalMiddlewares(adminRepository domain.AdminRepository, logger logger.ILogger, jwtFactory jwt.JwtServiceFactory) *ProviderMiddlewares {
	return &ProviderMiddlewares{
		adminRepository: adminRepository,
		logger:          logger,
		jwtFactory:      jwtFactory,
	}
}
