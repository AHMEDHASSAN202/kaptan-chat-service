package jwt

import (
	"samm/pkg/config"
	"samm/pkg/logger"
)

type JwtServiceFactory interface {
	AdminJwtService() JwtService
	PortalJwtService() JwtService
}

type JwtFactory struct {
	jWTConfig *config.JWTConfig
	logger    logger.ILogger
}

func (f *JwtFactory) AdminJwtService() JwtService {
	return &AdminJwtService{
		secretKey:    f.jWTConfig.AdminSigningKey,
		ExpiredHours: f.jWTConfig.AdminExpires,
		logger:       f.logger,
	}
}

func (f *JwtFactory) PortalJwtService() JwtService {
	return &PortalJwtService{
		secretKey:    f.jWTConfig.PortalSigningKey,
		ExpiredHours: f.jWTConfig.PortalExpires,
		logger:       f.logger,
	}
}

func NewJwtFactoryService(jWTConfig *config.JWTConfig, logger logger.ILogger) JwtServiceFactory {
	return &JwtFactory{
		jWTConfig: jWTConfig,
		logger:    logger,
	}
}
