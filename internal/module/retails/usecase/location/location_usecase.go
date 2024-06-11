package location

import (
	"context"
	"samm/internal/module/retails/consts"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type LocationUseCase struct {
	repo   domain.LocationRepository
	logger logger.ILogger
}

func (l LocationUseCase) StoreLocation(ctx context.Context, payload *location.StoreLocationDto) (err validators.ErrorResponse) {

	errRe := l.repo.StoreLocation(ctx, LocationBuilder(payload))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l LocationUseCase) BulkStoreLocation(ctx context.Context, payload []location.StoreLocationDto) (err validators.ErrorResponse) {

	data := make([]domain.Location, 0)
	for _, itemDoc := range payload {
		data = append(data, *LocationBuilder(&itemDoc))
	}

	errRe := l.repo.BulkStoreLocation(ctx, data)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l LocationUseCase) UpdateLocation(ctx context.Context, id string, payload *location.StoreLocationDto) (err validators.ErrorResponse) {
	domainLocation, errRe := l.repo.FindLocation(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}

	errRe = l.repo.UpdateLocation(ctx, UpdateLocationBuilder(payload, domainLocation))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}
func (l LocationUseCase) ToggleLocationStatus(ctx context.Context, id string) (err validators.ErrorResponse) {
	domainLocation, errRe := l.repo.FindLocation(ctx, utils.ConvertStringIdToObjectId(id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	if domainLocation.Status == consts.LocationStatusActive {
		domainLocation.Status = consts.LocationStatusInActive
	} else {
		domainLocation.Status = consts.LocationStatusActive
	}
	// todo handle admin_details for status change

	errRe = l.repo.UpdateLocation(ctx, domainLocation)
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return
}

func (l LocationUseCase) FindLocation(ctx context.Context, Id string) (location domain.Location, err validators.ErrorResponse) {
	domainLocation, errRe := l.repo.FindLocation(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return *domainLocation, validators.GetErrorResponseFromErr(errRe)
	}
	return *domainLocation, validators.ErrorResponse{}
}

func (l LocationUseCase) DeleteLocation(ctx context.Context, Id string) (err validators.ErrorResponse) {
	errRe := l.repo.DeleteLocation(ctx, utils.ConvertStringIdToObjectId(Id))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return validators.ErrorResponse{}
}

func (l LocationUseCase) DeleteLocationByAccountId(ctx context.Context, AccountId string) (err validators.ErrorResponse) {
	errRe := l.repo.DeleteLocationByAccountId(ctx, utils.ConvertStringIdToObjectId(AccountId))
	if errRe != nil {
		return validators.GetErrorResponseFromErr(errRe)
	}
	return validators.ErrorResponse{}
}

func (l LocationUseCase) ListLocation(ctx context.Context, payload *location.ListLocationDto) (locations []domain.Location, paginationResult utils.PaginationResult, err validators.ErrorResponse) {
	results, paginationResult, errRe := l.repo.ListLocation(ctx, payload)
	if errRe != nil {
		return results, paginationResult, validators.GetErrorResponseFromErr(errRe)
	}
	return results, paginationResult, validators.ErrorResponse{}

}
func (l *LocationUseCase) ToggleSnooze(ctx context.Context, dto *location.LocationToggleSnoozeDto) validators.ErrorResponse {
	locationDomain, err := l.repo.FindLocation(ctx, dto.Id)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	doc := domainBuilderToggleSnooze(dto, locationDomain)
	err = l.repo.UpdateLocation(ctx, doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

const tag = " LocationUseCase "

func NewLocationUseCase(repo domain.LocationRepository, logger logger.ILogger) domain.LocationUseCase {
	return &LocationUseCase{
		repo:   repo,
		logger: logger,
	}
}
