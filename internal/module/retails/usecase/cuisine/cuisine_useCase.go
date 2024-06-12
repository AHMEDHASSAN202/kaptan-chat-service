package cuisine

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strings"
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

func (oRec *CuisineUseCase) Create(ctx *context.Context, dto *cuisine.CreateCuisineDto) (*domain.Cuisine, validators.ErrorResponse) {
	doc := convertDtoArrToCorrespondingDomain(dto)
	err := oRec.repo.Create(doc)
	if err != nil {
		return doc, validators.GetErrorResponseFromErr(err)
	}
	return doc, validators.ErrorResponse{}
}

func (oRec *CuisineUseCase) Update(ctx *context.Context, dto *cuisine.UpdateCuisineDto) validators.ErrorResponse {
	findCuisine, findCuisineErr := oRec.GetById(ctx, dto.Id)
	if findCuisineErr.IsError {
		return findCuisineErr
	}
	doc := domainBuilderAtUpdate(dto, findCuisine)
	err := oRec.repo.Update(doc)
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

func (oRec *CuisineUseCase) List(ctx *context.Context, dto *cuisine.ListCuisinesDto) (cuisines *[]domain.Cuisine, paginationMeta *PaginationData, err validators.ErrorResponse) {
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

func (oRec *CuisineUseCase) CheckExists(ctx *context.Context, ids []string) validators.ErrorResponse {
	objIds := utils.ConvertStringIdsToObjectIds(ids)
	cuisines, err := oRec.repo.GetByIds(ctx, &objIds)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	if len(*cuisines) == 0 {
		return validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}
	cuisineIds := getCuisinesIds(cuisines)
	diffIds := utils.ElementsDiff(ids, cuisineIds)
	if diffIds != nil && len(diffIds) > 0 {
		oRec.logger.Error(strings.Join(diffIds, ", ") + " cuisine not exist")
		return validators.GetErrorResponseFromErr(errors.New(localization.E1002))
	}
	return validators.ErrorResponse{}
}
