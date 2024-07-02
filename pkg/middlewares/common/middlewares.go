package commmon

import (
	"samm/pkg/logger"
)

type ProviderMiddlewares struct {
	logger logger.ILogger
}

func NewCommonMiddlewares(logger logger.ILogger) *ProviderMiddlewares {
	return &ProviderMiddlewares{
		logger: logger,
	}
}
