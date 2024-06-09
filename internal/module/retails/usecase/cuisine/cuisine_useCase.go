package cuisine

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type CuisineUseCase struct {
	repo   domain.CuisineRepository
	logger logger.ILogger
}

func NewCuisineUseCase(repo domain.CuisineRepository, logger logger.ILogger) domain.CuisineUseCase {
	return &CuisineUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *CuisineUseCase) Create(ctx *context.Context, dto *[]cuisine.CreateCuisineDto) validators.ErrorResponse {
	err := oRec.repo.Create(ctx, convertDtoArrToCorrespondingDomain(dto))
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) Update(ctx *context.Context, dto *cuisine.UpdateCuisineDto) validators.ErrorResponse {
	findCuisine, findCuisineErr := oRec.GetById(ctx, dto.Id)
	if findCuisineErr.IsError {
		return findCuisineErr
	}
	doc := domainBuilderAtUpdate(dto, findCuisine)
	err := oRec.repo.Update(ctx, doc.ID, doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}
func (oRec *CuisineUseCase) SoftDelete(ctx *context.Context, id string) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	err := oRec.repo.SoftDelete(ctx, idDoc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) validators.ErrorResponse {
	err := oRec.repo.ChangeStatus(ctx, dto)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) List(ctx *context.Context, dto *cuisine.ListCuisinesDto) (cuisines *[]domain.Cuisine, paginationMeta *utils.PaginationResult, err validators.ErrorResponse) {
	cuisines, paginationMeta, resErr := oRec.repo.List(ctx, dto)
	if resErr != nil {
		err = validators.GetErrorResponseFromErr(resErr)
		return
	}
	return
}

func (oRec *CuisineUseCase) Find(ctx *context.Context, id string) (*domain.Cuisine, validators.ErrorResponse) {
	cuisine, err := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(id))
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	if cuisine == nil {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return cuisine, validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) GetById(ctx *context.Context, id string) (*domain.Cuisine, validators.ErrorResponse) {
	cuisines, err := oRec.repo.GetByIds(ctx, &[]primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	cuisinesVal := *cuisines
	if len(cuisinesVal) == 0 {
		return nil, validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}

	return &cuisinesVal[0], validators.ErrorResponse{}
}
