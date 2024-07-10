package kitchen

import (
	"context"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/internal/module/kitchen/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"time"
)

type KitchenUseCase struct {
	repo   domain.KitchenRepository
	logger logger.ILogger
}

const tag = " KitchenUseCase "

func NewKitchenUseCase(repo domain.KitchenRepository, logger logger.ILogger) domain.KitchenUseCase {
	return &KitchenUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (l KitchenUseCase) CreateKitchen(ctx context.Context, payload *kitchen.StoreKitchenDto) (err validators.ErrorResponse) {
	kitchenDomain := domain.Kitchen{}
	kitchenDomain.Name.Ar = payload.Name.Ar
	kitchenDomain.Name.En = payload.Name.En
	kitchenDomain.Email = payload.Email
	password, er := utils.HashPassword(payload.Password)
	if er != nil {
		return validators.GetErrorResponseFromErr(er)
	}
	kitchenDomain.Password = password
	kitchenDomain.CreatedAt = time.Now()
	kitchenDomain.UpdatedAt = time.Now()

	dbErr := l.repo.CreateKitchen(&kitchenDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}

func (l KitchenUseCase) UpdateKitchen(ctx context.Context, id string, payload *kitchen.UpdateKitchenDto) (err validators.ErrorResponse) {
	kitchenDomain, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(id))
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	kitchenDomain.Name.Ar = payload.Name.Ar
	kitchenDomain.Name.En = payload.Name.En
	kitchenDomain.Email = payload.Email

	if payload.Password != "" {
		password, er := utils.HashPassword(payload.Password)
		if er != nil {
			return validators.GetErrorResponseFromErr(er)
		}
		kitchenDomain.Password = password
	}
	kitchenDomain.UpdatedAt = time.Now()

	dbErr = l.repo.UpdateKitchen(kitchenDomain)
	if dbErr != nil {
		return validators.GetErrorResponseFromErr(dbErr)
	}
	return
}
func (l KitchenUseCase) FindKitchen(ctx context.Context, Id string) (kitchen domain.Kitchen, err validators.ErrorResponse) {
	domainKitchen, dbErr := l.repo.FindKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if dbErr != nil {
		return *domainKitchen, validators.GetErrorResponseFromErr(dbErr)
	}
	return *domainKitchen, validators.ErrorResponse{}
}

func (l KitchenUseCase) DeleteKitchen(ctx context.Context, Id string) (err validators.ErrorResponse) {

	delErr := l.repo.DeleteKitchen(ctx, utils.ConvertStringIdToObjectId(Id))
	if delErr != nil {
		return validators.GetErrorResponseFromErr(delErr)
	}
	return validators.ErrorResponse{}
}

func (l KitchenUseCase) List(ctx *context.Context, dto *kitchen.ListKitchenDto) (*responses.ListResponse, validators.ErrorResponse) {
	users, paginationMeta, resErr := l.repo.List(ctx, dto)
	if resErr != nil {
		return nil, validators.GetErrorResponseFromErr(resErr)
	}
	return responses.SetListResponse(users, paginationMeta), validators.ErrorResponse{}
}

