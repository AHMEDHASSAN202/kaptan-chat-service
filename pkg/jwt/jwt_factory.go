package jwt

import (
	"kaptan/pkg/config"
	"kaptan/pkg/logger"
)

type JwtServiceFactory interface {
	AdminJwtService() JwtService
	PortalJwtService() JwtService
	UserJwtService() JwtService
	KitchenJwtService() JwtService
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

func (f *JwtFactory) UserJwtService() JwtService {
	return &UserJwtService{
		secretKey:        f.jWTConfig.UserSigningKey,
		expiredHours:     f.jWTConfig.UserExpires,
		tempSecretKey:    f.jWTConfig.UserTempSigningKey,
		tempExpiredHours: f.jWTConfig.UserTempExpires,
		logger:           f.logger,
	}
}

func (f *JwtFactory) KitchenJwtService() JwtService {
	return &KitchenJwtService{
		secretKey:    f.jWTConfig.KitchenSigningKey,
		expiredHours: f.jWTConfig.PortalExpires,
		logger:       f.logger,
	}
}

func NewJwtFactoryService(jWTConfig *config.JWTConfig, logger logger.ILogger) JwtServiceFactory {
	return &JwtFactory{
		jWTConfig: jWTConfig,
		logger:    logger,
	}
}
