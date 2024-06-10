package common

import (
	"context"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

const tag = "CommonUseCase "

func NewCommonUseCase(repo domain.CommonRepository, logger logger.ILogger) domain.CommonUseCase {
	return &CommonUseCase{
		repo:   repo,
		logger: logger,
	}
}

type CommonUseCase struct {
	repo   domain.CommonRepository
	logger logger.ILogger
}

func (l CommonUseCase) ListCities(ctx context.Context, payload *location.ListCitiesDto) (data interface{}, err validators.ErrorResponse) {

	return CitiesBuilder(payload), validators.ErrorResponse{}

}
func (l CommonUseCase) ListCountries(ctx context.Context) (data interface{}, err validators.ErrorResponse) {

	return CountriesBuilder(), validators.ErrorResponse{}

}
